// Copyright 2020 PingCAP, Inc. Licensed under MIT.

// https://github.com/grafana/grafana/blob/v6.1.6/public/app/features/dashboard/state/PanelModel.ts
declare class PanelModel {
    snapshotData?: any
}

declare interface Annotation {
    snapshotData?: any
}

declare interface Annotations {
    list: Annotation[]
}

// https://github.com/grafana/grafana/blob/v6.1.6/public/app/features/dashboard/state/DashboardModel.ts
declare class DashboardModel {
    readonly refresh?: string
    title: string
    id: number | null
    uid: string | null
    readonly annotations: Annotations
    time: TimeRange
    snapshot?: Snapshot

    getSaveModelClone(): DashboardModel
    expandRows(): void
    startRefresh(): void
    forEachPanel(callback: (panel: PanelModel, index: number) => void): void
    removePanel(_: null): void
}

declare interface TimeRange {}

// https://github.com/grafana/grafana/blob/v6.1.6/public/app/features/dashboard/services/TimeSrv.ts
declare class TimeSrv {
    readonly dashboard: DashboardModel

    setAutoRefresh(interval?: string): void
    timeRange(): TimeRange
}

declare interface JQCoords {
    left: number
    top: number
}

declare interface JQ extends JQLite {
    filter(selector: string): JQ
    hide(): JQ
    show(): JQ
}

// https://github.com/grafana/grafana/blob/v6.1.6/public/app/features/dashboard/components/ShareModal/ShareSnapshotCtrl.ts#L38
declare interface Snapshot {
    readonly timestamp: Date
}
