import React from 'react'

import { cleanup, fireEvent } from '@testing-library/react'

import { renderWithBrandedContext } from '@sourcegraph/shared/src/testing'

import { Panel } from './Panel'
import { panels, panelProps } from './Panel.fixtures'

describe('Panel', () => {
    const location = {
        pathname: `/${panelProps.repoName}`,
        search: '?L4:7',
        hash: `#tab=${panels[0].id}`,
    }
    const route = `${location.pathname}${location.search}${location.hash}`

    afterEach(cleanup)

    it('preserves `location.pathname` and `location.hash` on tab change', async () => {
        const renderResult = renderWithBrandedContext(<Panel {...panelProps} />, { route })

        const panelToSelect = panels[2]
        const panelButton = await renderResult.findByRole('tab', { name: panelToSelect.title })
        fireEvent.click(panelButton)

        expect(renderResult.history.location.pathname).toEqual(location.pathname)
        expect(renderResult.history.location.search).toEqual(location.search)
        expect(renderResult.history.location.hash).toEqual(`#tab=${panelToSelect.id}`)
    })
})
