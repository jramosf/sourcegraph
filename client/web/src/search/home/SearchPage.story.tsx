import { storiesOf } from '@storybook/react'
import { parseISO } from 'date-fns'
import { createMemoryHistory } from 'history'
import React from 'react'
import { getDocumentNode } from '@sourcegraph/http-client'
import { MockedTestProvider } from '@sourcegraph/shared/src/testing/apollo'

import { NOOP_TELEMETRY_SERVICE } from '@sourcegraph/shared/src/telemetry/telemetryService'
import {
    mockFetchAutoDefinedSearchContexts,
    mockFetchSearchContexts,
    mockGetUserSearchContextNamespaces,
} from '@sourcegraph/shared/src/testing/searchContexts/testHelpers'
import { extensionsController } from '@sourcegraph/shared/src/testing/searchTestHelpers'
import { ThemeProps } from '@sourcegraph/shared/src/theme'
import {
    HOME_PANELS_QUERY,
    RECENTLY_SEARCHED_REPOSITORIES_TO_LOAD,
    RECENT_SEARCHES_TO_LOAD,
} from '../../search/panels/HomePanels'
import { WebStory } from '../../components/WebStory'
import { FeatureFlagName } from '../../featureFlags/featureFlags'
import { SourcegraphContext } from '../../jscontext'
import { useExperimentalFeatures } from '../../stores'
import { ThemePreference } from '../../stores/themeState'
import {
    _fetchRecentFileViews,
    _fetchSavedSearches,
    _fetchCollaborators,
    authUser,
    recentSearchesPayload,
} from '../panels/utils'

import { SearchPage, SearchPageProps } from './SearchPage'

const history = createMemoryHistory()
const defaultProps = (props: ThemeProps): SearchPageProps => ({
    isSourcegraphDotCom: false,
    settingsCascade: {
        final: null,
        subjects: null,
    },
    location: history.location,
    history,
    extensionsController,
    telemetryService: NOOP_TELEMETRY_SERVICE,
    themePreference: ThemePreference.Light,
    onThemePreferenceChange: () => undefined,
    authenticatedUser: authUser,
    globbing: false,
    platformContext: {} as any,
    keyboardShortcuts: [],
    searchContextsEnabled: true,
    selectedSearchContextSpec: '',
    setSelectedSearchContextSpec: () => {},
    defaultSearchContextSpec: '',
    isLightTheme: props.isLightTheme,
    fetchSavedSearches: _fetchSavedSearches,
    fetchRecentFileViews: _fetchRecentFileViews,
    fetchCollaborators: _fetchCollaborators,
    now: () => parseISO('2020-09-16T23:15:01Z'),
    fetchAutoDefinedSearchContexts: mockFetchAutoDefinedSearchContexts(),
    fetchSearchContexts: mockFetchSearchContexts,
    hasUserAddedRepositories: false,
    hasUserAddedExternalServices: false,
    getUserSearchContextNamespaces: mockGetUserSearchContextNamespaces,
    featureFlags: new Map<FeatureFlagName, boolean>(),
})

if (!window.context) {
    // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
    window.context = {} as SourcegraphContext & Mocha.SuiteFunction
}
window.context.allowSignup = true

const { add } = storiesOf('web/search/home/SearchPage', module)
    .addParameters({
        design: {
            type: 'figma',
            url: 'https://www.figma.com/file/sPRyyv3nt5h0284nqEuAXE/12192-Sourcegraph-server-page-v1?node-id=255%3A3',
        },
        chromatic: { viewports: [544, 577, 769, 993], disableSnapshot: false },
    })
    .addDecorator(Story => {
        useExperimentalFeatures.setState({ showSearchContext: false, showEnterpriseHomePanels: false })
        return <Story />
    })

const mocks = [
    {
        request: {
            query: getDocumentNode(HOME_PANELS_QUERY),
            variables: {
                userId: '0',
                firstRecentlySearchedRepositories: RECENTLY_SEARCHED_REPOSITORIES_TO_LOAD,
                firstRecentSearches: RECENT_SEARCHES_TO_LOAD,
            },
        },
        result: {
            data: {
                node: {
                    __typename: 'User',
                    recentlySearchedRepositoriesLogs: recentSearchesPayload(),
                    recentSearchesLogs: recentSearchesPayload(),
                },
            },
        },
    },
]

add('Cloud with panels', () => (
    <WebStory>
        {webProps => {
            useExperimentalFeatures.setState({ showEnterpriseHomePanels: true })
            return (
                <MockedTestProvider mocks={mocks}>
                    <SearchPage {...defaultProps(webProps)} isSourcegraphDotCom={true} />
                </MockedTestProvider>
            )
        }}
    </WebStory>
))

add('Cloud with panels and collaborators', () => (
    <WebStory>
        {webProps => {
            useExperimentalFeatures.setState({ showEnterpriseHomePanels: true })
            useExperimentalFeatures.setState({ homepageUserInvitation: true })
            return (
                <MockedTestProvider mocks={mocks}>
                    <SearchPage {...defaultProps(webProps)} isSourcegraphDotCom={true} />
                </MockedTestProvider>
            )
        }}
    </WebStory>
))

add('Cloud marketing home', () => (
    <WebStory>
        {webProps => <SearchPage {...defaultProps(webProps)} isSourcegraphDotCom={true} authenticatedUser={null} />}
    </WebStory>
))

add('Server with panels', () => (
    <WebStory>
        {webProps => {
            useExperimentalFeatures.setState({ showEnterpriseHomePanels: true })
            return (
                <MockedTestProvider mocks={mocks}>
                    <SearchPage {...defaultProps(webProps)} />
                </MockedTestProvider>
            )
        }}
    </WebStory>
))
