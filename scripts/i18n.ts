// Copyright 2020 PingCAP, Inc. Licensed under MIT.

interface Label {
    readonly unsupported: string
    readonly exportTitle: string
    readonly exportSubtitle: string
    readonly expandAllPanels: string
    readonly numberOfPanels: string
    readonly loadingPanels: string
    readonly export: string
    readonly exportImmediately: string
    readonly cancel: string
}

type Language = 'en-US' | 'zh-TW' | 'zh-CN';

const LABELS: {[_ in Language]: Label} = {
    'en-US': {
        unsupported: 'Script failed to run. Please ensure it is running on the Grafana v6.x.x dashboard page.',
        exportTitle: 'Going to export snapshot of the current dashboard: ',
        exportSubtitle: 'Note: Only metrics from visible panels can be exported',
        numberOfPanels: 'Number of panels',
        loadingPanels: 'Loading panels',
        expandAllPanels: 'Expand all rows',
        export: 'Export',
        exportImmediately: 'Export immediately',
        cancel: 'Cancel',
    },
    'zh-TW': {
        unsupported: '執行失敗。請確保是在 Grafana v6.x.x 的儀表板上執行劇本。',
        exportTitle: '即將匯出當前儀表板：',
        exportSubtitle: '注意：只能匯出已展開的面板內的資料',
        numberOfPanels: '面板數量',
        loadingPanels: '載入面板中',
        expandAllPanels: '展開所有面板',
        export: '匯出',
        exportImmediately: '立即匯出',
        cancel: '取消',
    },
    'zh-CN': {
        unsupported: '执行失败。请确保是在 Grafana v6.x.x 的监控网页上执行脚本。',
        exportTitle: '将导出当前的监控页面：',
        exportSubtitle: '注意：只有已展开面板内的监控会被导出',
        numberOfPanels: '面板数量',
        loadingPanels: '加载面板中',
        expandAllPanels: '展开所有面板',
        export: '导出',
        exportImmediately: '立刻导出',
        cancel: '取消',
    },
};

function getLanguage(): Language {
    const browserLanguage = navigator.language;
    if (/^zh-/.test(browserLanguage)) {
        if (/\b(TW|HK|MO|Hant)\b/.test(browserLanguage)) {
            return 'zh-TW';
        }
        return 'zh-CN';
    }
    return 'en-US';
}

const L = LABELS[getLanguage()];
