import isEqual from 'lodash/isEqual'
import { Observable } from 'rxjs/Observable'
import { distinctUntilChanged } from 'rxjs/operators/distinctUntilChanged'
import { map } from 'rxjs/operators/map'
import { mergeMap } from 'rxjs/operators/mergeMap'
import { gql, queryGraphQL } from '../backend/graphql'
import { mutateConfigurationGraphQL } from '../configuration/backend'
import { currentConfiguration, SavedQueryConfiguration } from '../settings/configuration'
import { SearchOptions } from './index'

export function searchText(options: SearchOptions): Observable<GQL.ISearchResults> {
    return queryGraphQL(
        gql`
            query Search($query: String!) {
                search(query: $query) {
                    results {
                        limitHit
                        missing
                        cloning
                        timedout
                        results {
                            __typename
                            ... on FileMatch {
                                resource
                                limitHit
                                lineMatches {
                                    preview
                                    lineNumber
                                    offsetAndLengths
                                }
                            }
                            ... on CommitSearchResult {
                                refs {
                                    name
                                    displayName
                                    prefix
                                    repository {
                                        uri
                                    }
                                }
                                sourceRefs {
                                    name
                                    displayName
                                    prefix
                                    repository {
                                        uri
                                    }
                                }
                                messagePreview {
                                    value
                                    highlights {
                                        line
                                        character
                                        length
                                    }
                                }
                                diffPreview {
                                    value
                                    highlights {
                                        line
                                        character
                                        length
                                    }
                                }
                                commit {
                                    repository {
                                        uri
                                    }
                                    oid
                                    abbreviatedOID
                                    author {
                                        person {
                                            displayName
                                            avatarURL
                                        }
                                        date
                                    }
                                    message
                                }
                            }
                        }
                        alert {
                            title
                            description
                            proposedQueries {
                                description
                                query {
                                    query
                                }
                            }
                        }
                    }
                }
            }
        `,
        { query: options.query }
    ).pipe(
        map(({ data, errors }) => {
            if (!data || !data.search || !data.search.results) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.search.results
        })
    )
}

export function fetchSearchResultStats(options: SearchOptions): Observable<GQL.ISearchResults> {
    return queryGraphQL(
        gql`
            query SearchResultsCount($query: String!) {
                search(query: $query) {
                    results {
                        limitHit
                        missing
                        cloning
                        resultCount
                        approximateResultCount
                        sparkline
                    }
                }
            }
        `,
        { query: options.query }
    ).pipe(
        map(({ data, errors }) => {
            if (!data || !data.search || !data.search.results) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.search.results
        })
    )
}

export function fetchSuggestions(options: SearchOptions): Observable<GQL.SearchSuggestion> {
    return queryGraphQL(
        gql`
            query Search($query: String!) {
                search(query: $query) {
                    suggestions {
                        ... on Repository {
                            __typename
                            uri
                        }
                        ... on File {
                            __typename
                            name
                            isDirectory
                            repository {
                                uri
                            }
                        }
                    }
                }
            }
        `,
        { query: options.query }
    ).pipe(
        mergeMap(({ data, errors }) => {
            if (!data || !data.search || !data.search.suggestions) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.search.suggestions
        })
    )
}

export function fetchSearchScopes(): Observable<GQL.ISearchScope[]> {
    return queryGraphQL(gql`
        query SearchScopes {
            searchScopes {
                name
                value
            }
        }
    `).pipe(
        map(({ data, errors }) => {
            if (!data || !data.searchScopes) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.searchScopes
        })
    )
}

export function fetchRepoGroups(): Observable<GQL.IRepoGroup[]> {
    return queryGraphQL(gql`
        query RepoGroups {
            repoGroups {
                name
                repositories
            }
        }
    `).pipe(
        map(({ data, errors }) => {
            if (!data || !data.repoGroups) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.repoGroups
        })
    )
}

const savedQueryFragment = gql`
    fragment SavedQueryFields on SavedQuery {
        id
        subject {
            ... on Org {
                id
            }
            ... on User {
                id
            }
        }
        index
        description
        showOnHomepage
        query {
            query
        }
    }
`

function savedQueriesEqual(a: SavedQueryConfiguration, b: SavedQueryConfiguration): boolean {
    return isEqual(a, b)
}

export function observeSavedQueries(): Observable<GQL.ISavedQuery[]> {
    return currentConfiguration.pipe(
        map(config => config['search.savedQueries']),
        distinctUntilChanged(
            (a, b) =>
                (!a && !b) || (!!a && !!b && a.length === b.length && a.every((q, i) => savedQueriesEqual(q, b[i])))
        ),
        mergeMap(fetchSavedQueries)
    )
}

function fetchSavedQueries(): Observable<GQL.ISavedQuery[]> {
    return queryGraphQL(gql`
        query SavedQueries {
            savedQueries {
                ...SavedQueryFields
            }
        }
        ${savedQueryFragment}
    `).pipe(
        map(({ data, errors }) => {
            if (!data || !data.savedQueries) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.savedQueries
        })
    )
}

export function createSavedQuery(
    subject: GQL.ConfigurationSubject | GQL.IConfigurationSubject | { id: GQLID },
    description: string,
    query: string,
    showOnHomepage: boolean
): Observable<GQL.ISavedQuery> {
    return mutateConfigurationGraphQL(
        subject,
        gql`
            mutation CreateSavedQuery(
                $subject: ID!
                $lastID: Int
                $description: String!
                $query: String!
                $showOnHomepage: Boolean
            ) {
                configurationMutation(input: { subject: $subject, lastID: $lastID }) {
                    createSavedQuery(description: $description, query: $query, showOnHomepage: $showOnHomepage) {
                        ...SavedQueryFields
                    }
                }
            }
            ${savedQueryFragment}
        `,
        { description, query, showOnHomepage }
    ).pipe(
        map(({ data, errors }) => {
            if (!data || !data.configurationMutation || !data.configurationMutation.createSavedQuery) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.configurationMutation.createSavedQuery
        })
    )
}

export function updateSavedQuery(
    subject: GQL.ConfigurationSubject | GQL.IConfigurationSubject | { id: GQLID },
    id: GQLID,
    description: string,
    query: string,
    showOnHomepage: boolean
): Observable<GQL.ISavedQuery> {
    return mutateConfigurationGraphQL(
        subject,
        gql`
            mutation UpdateSavedQuery(
                $subject: ID!
                $lastID: Int
                $id: ID!
                $description: String
                $query: String
                $showOnHomepage: Boolean
            ) {
                configurationMutation(input: { subject: $subject, lastID: $lastID }) {
                    updateSavedQuery(
                        id: $id
                        description: $description
                        query: $query
                        showOnHomepage: $showOnHomepage
                    ) {
                        ...SavedQueryFields
                    }
                }
            }
            ${savedQueryFragment}
        `,
        { id, description, query, showOnHomepage }
    ).pipe(
        map(({ data, errors }) => {
            if (!data || !data.configurationMutation || !data.configurationMutation.updateSavedQuery) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
            return data.configurationMutation.updateSavedQuery
        })
    )
}

export function deleteSavedQuery(
    subject: GQL.ConfigurationSubject | GQL.IConfigurationSubject | { id: GQLID },
    id: GQLID
): Observable<void> {
    return mutateConfigurationGraphQL(
        subject,
        gql`
            mutation DeleteSavedQuery($subject: ID!, $lastID: Int, $id: ID!) {
                configurationMutation(input: { subject: $subject, lastID: $lastID }) {
                    deleteSavedQuery(id: $id) {
                        alwaysNil
                    }
                }
            }
        `,
        { id }
    ).pipe(
        map(({ data, errors }) => {
            if (!data || !data.configurationMutation || !data.configurationMutation.deleteSavedQuery) {
                throw Object.assign(new Error((errors || []).map(e => e.message).join('\n')), { errors })
            }
        })
    )
}
