import classNames from 'classnames'
import AddIcon from 'mdi-react/AddIcon'
import CloseIcon from 'mdi-react/CloseIcon'
import EmailOpenOutlineIcon from 'mdi-react/EmailOpenOutlineIcon'
import React, { useCallback, useState } from 'react'
import { map } from 'rxjs/operators'

import { ErrorAlert } from '@sourcegraph/branded/src/components/alerts'
import { Form } from '@sourcegraph/branded/src/components/Form'
import { asError, createAggregateError, isErrorLike } from '@sourcegraph/common'
import { Link } from '@sourcegraph/shared/src/components/Link'
import { Scalars } from '@sourcegraph/shared/src/graphql-operations'
import { gql } from '@sourcegraph/shared/src/graphql/graphql'
import { Button, LoadingSpinner } from '@sourcegraph/wildcard'

import { AuthenticatedUser } from '../../../auth'
import { requestGraphQL } from '../../../backend/graphql'
import { CopyableText } from '../../../components/CopyableText'
import { DismissibleAlert } from '../../../components/DismissibleAlert'
import {
    InviteUserToOrganizationResult,
    InviteUserToOrganizationVariables,
    AddUserToOrganizationResult,
    AddUserToOrganizationVariables,
    InviteUserToOrganizationFields,
} from '../../../graphql-operations'
import { eventLogger } from '../../../tracking/eventLogger'

import styles from './InviteForm.module.scss'

const emailInvitesEnabled = window.context.emailEnabled

interface Invited extends InviteUserToOrganizationFields {
    username: string
}

interface Props {
    orgID: Scalars['ID']
    authenticatedUser: AuthenticatedUser | null

    /** Called when the organization members list changes. */
    onDidUpdateOrganizationMembers: () => void

    onOrganizationUpdate: () => void
}

export const InviteForm: React.FunctionComponent<Props> = ({
    orgID,
    authenticatedUser,
    onDidUpdateOrganizationMembers,
    onOrganizationUpdate,
}) => {
    const [username, setUsername] = useState<string>('')
    const onUsernameChange = useCallback<React.ChangeEventHandler<HTMLInputElement>>(event => {
        setUsername(event.currentTarget.value)
    }, [])
    const [loading, setLoading] = useState<'inviteUserToOrganization' | 'addUserToOrganization' | Error>()
    const [invited, setInvited] = useState<Invited>()
    const [isInviteShown, setShowInvitation] = useState<boolean>(false)

    const inviteUser = useCallback(() => {
        eventLogger.log('InviteOrgMemberClicked')
        ;(async () => {
            setLoading('inviteUserToOrganization')
            const { invitationURL, sentInvitationEmail } = await inviteUserToOrganization(username, orgID)
            setInvited({ username, sentInvitationEmail, invitationURL })
            setShowInvitation(true)
            onOrganizationUpdate()
            setUsername('')
            setLoading(undefined)
        })().catch(error => setLoading(asError(error)))
    }, [onOrganizationUpdate, orgID, setShowInvitation, username])

    const onInviteClick = useCallback<React.MouseEventHandler<HTMLButtonElement>>(() => {
        inviteUser()
    }, [inviteUser])

    const dismissNotification = (): void => {
        setShowInvitation(false)
    }

    const viewerCanAddUserToOrganization = !!authenticatedUser && authenticatedUser.siteAdmin

    const onSubmit = useCallback<React.FormEventHandler<HTMLFormElement>>(
        async event => {
            event.preventDefault()
            if (!viewerCanAddUserToOrganization) {
                inviteUser()
            } else {
                setLoading('addUserToOrganization')
                try {
                    await addUserToOrganization(username, orgID)
                    onDidUpdateOrganizationMembers()
                    setUsername('')
                    setLoading(undefined)
                } catch (error) {
                    setLoading(asError(error))
                }
            }
        },
        [inviteUser, onDidUpdateOrganizationMembers, orgID, username, viewerCanAddUserToOrganization]
    )

    return (
        <div>
            <div className={styles.container}>
                <label htmlFor="invite-form__username">
                    {viewerCanAddUserToOrganization ? 'Add or invite member' : 'Invite member'}
                </label>
                <Form className="form-inline align-items-start" onSubmit={onSubmit}>
                    <input
                        type="text"
                        className="form-control mb-2 mr-sm-2"
                        id="invite-form__username"
                        placeholder="Username"
                        onChange={onUsernameChange}
                        value={username}
                        autoComplete="off"
                        autoCapitalize="off"
                        autoCorrect="off"
                        required={true}
                        spellCheck={false}
                        size={30}
                    />
                    <div className="d-block d-md-inline">
                        {viewerCanAddUserToOrganization && (
                            <Button
                                type="submit"
                                disabled={loading === 'addUserToOrganization' || loading === 'inviteUserToOrganization'}
                                className="mr-2"
                                data-tooltip="Add immediately without sending invitation (site admins only)"
                                variant="primary"
                            >
                                {loading === 'addUserToOrganization' ? (
                                    <LoadingSpinner />
                                ) : (
                                    <AddIcon className="icon-inline" />
                                )}{' '}
                                Add member
                            </Button>
                        )}
                        {(emailInvitesEnabled || !viewerCanAddUserToOrganization) && (
                            <Button
                                type={viewerCanAddUserToOrganization ? 'button' : 'submit'}
                                disabled={loading === 'addUserToOrganization' || loading === 'inviteUserToOrganization'}
                                className={viewerCanAddUserToOrganization ? 'btn-secondary' : 'btn-primary'}
                                data-tooltip={
                                    emailInvitesEnabled
                                        ? 'Send invitation email with link to join this organization'
                                        : 'Generate invitation link to manually send to user'
                                }
                                onClick={viewerCanAddUserToOrganization ? onInviteClick : undefined}
                            >
                                {loading === 'inviteUserToOrganization' ? (
                                    <LoadingSpinner />
                                ) : (
                                    <EmailOpenOutlineIcon className="icon-inline" />
                                )}{' '}
                                {emailInvitesEnabled
                                    ? viewerCanAddUserToOrganization
                                        ? 'Send invitation to join'
                                        : 'Send invitation'
                                    : 'Generate invitation link'}
                            </Button>
                        )}
                    </div>
                </Form>
            </div>
            {authenticatedUser?.siteAdmin && !emailInvitesEnabled && (
                <DismissibleAlert className="alert-info" partialStorageKey="org-invite-email-config">
                    <p className=" mb-0">
                        Set <code>email.smtp</code> in <Link to="/site-admin/configuration">site configuration</Link> to
                        send email notifications about invitations.
                    </p>
                </DismissibleAlert>
            )}
            {invited && isInviteShown && (
                <InvitedNotification
                    key={invited.username}
                    {...invited}
                    className={classNames('alert alert-success', styles.alert)}
                    onDismiss={dismissNotification}
                />
            )}
            {isErrorLike(loading) && <ErrorAlert className={styles.alert} error={loading} />}
        </div>
    )
}

function inviteUserToOrganization(
    username: string,
    organization: Scalars['ID']
): Promise<InviteUserToOrganizationResult['inviteUserToOrganization']> {
    return requestGraphQL<InviteUserToOrganizationResult, InviteUserToOrganizationVariables>(
        gql`
            mutation InviteUserToOrganization($organization: ID!, $username: String!) {
                inviteUserToOrganization(organization: $organization, username: $username) {
                    ...InviteUserToOrganizationFields
                }
            }

            fragment InviteUserToOrganizationFields on InviteUserToOrganizationResult {
                sentInvitationEmail
                invitationURL
            }
        `,
        {
            username,
            organization,
        }
    )
        .pipe(
            map(({ data, errors }) => {
                if (!data || !data.inviteUserToOrganization || (errors && errors.length > 0)) {
                    eventLogger.log('InviteOrgMemberFailed')
                    throw createAggregateError(errors)
                }
                eventLogger.log('OrgMemberInvited')
                return data.inviteUserToOrganization
            })
        )
        .toPromise()
}

function addUserToOrganization(username: string, organization: Scalars['ID']): Promise<void> {
    return requestGraphQL<AddUserToOrganizationResult, AddUserToOrganizationVariables>(
        gql`
            mutation AddUserToOrganization($organization: ID!, $username: String!) {
                addUserToOrganization(organization: $organization, username: $username) {
                    alwaysNil
                }
            }
        `,
        {
            username,
            organization,
        }
    )
        .pipe(
            map(({ data, errors }) => {
                if (!data || !data.addUserToOrganization || (errors && errors.length > 0)) {
                    eventLogger.log('AddOrgMemberFailed')
                    throw createAggregateError(errors)
                }
                eventLogger.log('OrgMemberAdded')
            })
        )
        .toPromise()
}

interface InvitedNotificationProps {
    username: string
    sentInvitationEmail: boolean
    invitationURL: string
    onDismiss: () => void
    className?: string
}

const InvitedNotification: React.FunctionComponent<InvitedNotificationProps> = ({
    className,
    username,
    sentInvitationEmail,
    invitationURL,
    onDismiss,
}) => (
    <div className={classNames(styles.invitedNotification, className)}>
        <div className={styles.message}>
            {sentInvitationEmail ? (
                <>
                    Invitation sent to {username}. You can also send {username} the invitation link directly:
                </>
            ) : (
                <>Generated invitation link. Copy and send it to {username}:</>
            )}
            <CopyableText text={invitationURL} size={40} className="mt-2" />
        </div>
        <Button className="btn-icon" title="Dismiss" onClick={onDismiss}>
            <CloseIcon className="icon-inline" />
        </Button>
    </div>
)
