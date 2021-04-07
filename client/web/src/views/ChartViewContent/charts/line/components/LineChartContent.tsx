import React, { ReactElement, useCallback, useMemo, useState, MouseEvent, useEffect } from 'react'
import classnames from 'classnames'
import { LineChartContent as LineChartContentType } from 'sourcegraph'
import { useThrottledCallback } from 'use-debounce'
import { curveLinear } from '@visx/curve'
import { RenderTooltipParams } from '@visx/xychart/lib/components/Tooltip'
import { Axis, GlyphSeries, LineSeries, Tooltip, XYChart } from '@visx/xychart'
import { GridScale } from '@visx/grid/lib/types'
import { format } from 'd3-format'
import { timeFormat } from 'd3-time-format'
import { GridColumns, GridRows } from '@visx/grid'
import { Group } from '@visx/group'
import { GlyphDot as Glyph } from '@visx/glyph'
import isValidNumber from '@visx/xychart/lib/typeguards/isValidNumber'
import { EventHandlerParams } from '@visx/xychart/lib/types'

import { generateAccessors } from '../helpers/generate-accessors'
import { GlyphDot } from './GlyphDot'
import { TooltipContent } from './TooltipContent'
import { onDatumClick } from '../../types'
import { DEFAULT_LINE_STROKE } from '../colors'
import { useScales } from '../helpers/use-scales'

// Chart configuration
const WIDTH_PER_TICK = 70
const HEIGHT_PER_TICK = 40
const MARGIN = { top: 10, left: 30, bottom: 20, right: 20 }
const SCALES_CONFIG = {
    x: {
        type: 'time' as const,
        nice: true,
    },
    y: {
        type: 'linear' as const,
        nice: true,
        zero: false,
        clamp: true,
    },
}

// Formatters
const dateFormatter = timeFormat('%d %b')
const formatDate = (date: Date): string => dateFormatter(date)

export interface LineChartContentProps<Datum extends object>
    extends Omit<LineChartContentType<Datum, keyof Datum>, 'chart'> {
    width: number
    height: number
    onDatumClick: onDatumClick
}

interface ActiveDatum<Datum extends object> extends EventHandlerParams<Datum> {
    line: LineChartContentProps<Datum>['series'][number]
}

export function LineChartContent<Datum extends object>(props: LineChartContentProps<Datum>): ReactElement {
    const { width, height, data, series, xAxis, onDatumClick } = props

    // Derived
    const innerWidth = width - MARGIN.left - MARGIN.right
    const innerHeight = height - MARGIN.top - MARGIN.bottom

    const numberOfTicksX = Math.max(1, Math.floor(innerWidth / WIDTH_PER_TICK))
    const numberOfTicksY = Math.max(1, Math.floor(innerHeight / HEIGHT_PER_TICK))

    const sortedData = useMemo(
        () => data.sort((firstDatum, secondDatum) => +firstDatum[xAxis.dataKey] - +secondDatum[xAxis.dataKey]),
        [data, xAxis]
    )
    const accessors = useMemo(() => generateAccessors(xAxis, series), [xAxis, series])

    const { config: scalesConfig, xScale, yScale } = useScales({
        config: SCALES_CONFIG,
        data: sortedData,
        width: innerWidth,
        height: innerHeight,
        accessors,
    })

    // state
    const [activeDatum, setActiveDatum] = useState<ActiveDatum<Datum> | null>(null)

    // callbacks
    const renderTooltip = useCallback(
        (renderProps: RenderTooltipParams<Datum>) => (
            <TooltipContent
                {...renderProps}
                accessors={accessors}
                series={series}
                className="line-chart__tooltip-content"
            />
        ),
        [accessors, series]
    )

    const handlePointerMove = useThrottledCallback(
        (event: EventHandlerParams<Datum>) => {
            const line = series.find(line => line.dataKey === event.key)

            if (!line) {
                return
            }

            setActiveDatum({
                ...event,
                line,
            })
        },
        100,
        { leading: true }
    )

    // Cancel scheduled handlePointerMove callback in case if component was removed
    // from render tree to avoid can't perform state update on an unmounted component error.
    useEffect(() => () => handlePointerMove.cancel(), [handlePointerMove])

    const handlePointerOut = useCallback(() => setActiveDatum(null), [setActiveDatum])
    const handlePointerUp = useCallback((info: EventHandlerParams<Datum>) => {
        info.event?.persist()

        // By types from visx/xychart index can be undefined
        const activeDatumIndex = activeDatum?.index
        const line = series.find(line => line.dataKey === info.key)

        if (!info.event || !line || !isValidNumber(activeDatumIndex)) {
            return
        }

        onDatumClick({
            originEvent: info.event as MouseEvent<unknown>,
            link: line?.linkURLs?.[activeDatumIndex],
        })
    }, [series, onDatumClick, activeDatum])

    const activeDatumLink = activeDatum?.line?.linkURLs?.[activeDatum?.index]
    const rootClasses = classnames('line-chart__content', { 'line-chart__content--with-cursor': !!activeDatumLink })

    return (
        <div className={rootClasses}>
            <XYChart
                xScale={scalesConfig.x}
                yScale={scalesConfig.y}
                height={height}
                width={width}
                captureEvents={true}
                margin={MARGIN}
                onPointerMove={handlePointerMove}
                onPointerUp={handlePointerUp}
                onPointerOut={handlePointerOut}
            >
                <Group top={MARGIN.top} left={MARGIN.left}>
                    <GridRows
                        scale={yScale as GridScale}
                        numTicks={numberOfTicksY}
                        width={innerWidth}
                        className="line-chart__grid-line"
                    />

                    <GridColumns
                        scale={xScale as GridScale}
                        numTicks={numberOfTicksX}
                        height={innerHeight}
                        className="line-chart__grid-line"
                    />
                </Group>

                <Axis
                    orientation="bottom"
                    tickValues={xScale.ticks(numberOfTicksX)}
                    tickFormat={formatDate}
                    numTicks={numberOfTicksX}
                    axisClassName="line-chart__axis"
                    axisLineClassName="line-chart__axis-line"
                    tickClassName="line-chart__axis-tick"
                />
                <Axis
                    orientation="left"
                    numTicks={numberOfTicksY}
                    tickFormat={format('~s')}
                    axisClassName="line-chart__axis"
                    axisLineClassName="line-chart__axis-line"
                    tickClassName="line-chart__axis-tick"
                />

                {series.map(line => (
                    <Group key={line.dataKey as string}>
                        <LineSeries
                            dataKey={line.dataKey as string}
                            data={sortedData}
                            strokeWidth={3}
                            xAccessor={accessors.x}
                            yAccessor={accessors.y[line.dataKey as string]}
                            stroke={line.stroke ?? DEFAULT_LINE_STROKE}
                            curve={curveLinear}
                        />

                        <GlyphSeries
                            dataKey={line.dataKey as string}
                            data={sortedData}
                            /* eslint-disable-next-line react/jsx-no-bind */
                            colorAccessor={() => line.stroke ?? DEFAULT_LINE_STROKE}
                            xAccessor={accessors.x}
                            yAccessor={accessors.y[line.dataKey as string]}
                            renderGlyph={GlyphDot}
                        />
                    </Group>
                ))}

                <Group top={MARGIN.top} left={MARGIN.left}>
                    {activeDatum && (
                        <Glyph
                            className="line-chart__glyph line-chart__glyph--active"
                            r={8}
                            fill={activeDatum.line.stroke ?? DEFAULT_LINE_STROKE}
                            cx={xScale(accessors.x(activeDatum.datum))}
                            cy={yScale(accessors.y[activeDatum.key](activeDatum.datum))}
                        />
                    )}
                </Group>

                <Tooltip
                    className="line-chart__tooltip"
                    showHorizontalCrosshair={false}
                    showVerticalCrosshair={true}
                    snapTooltipToDatumX={false}
                    snapTooltipToDatumY={false}
                    showDatumGlyph={false}
                    showSeriesGlyphs={false}
                    renderTooltip={renderTooltip}
                />
            </XYChart>
        </div>
    )
}
