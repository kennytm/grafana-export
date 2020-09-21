// Copyright 2020 PingCAP, Inc. Licensed under MIT.

// `import angular from 'angular'` can't be used due to --module none.
// the following tricks tsc into recognizing declarations of '@types/angular'
// without any real effect.
false && import('angular');

const enum ExportState {
    idle,
    started,
    shortcut,
}

try {
    const $ = angular.element;
    const FIRST_STYLESHEET = document.styleSheets[0];

    const STYLES = [
        '.__dexport_hint{z-index:2222;position:fixed;background:#fff;color:#000;padding:1em;font-size:18px;border-left:3px solid #faad14;right:0;top:5em;opacity:.8;min-width:30em}',
        '.__dexport_hint:hover{opacity:1}',
        '.__dexport_hint button{margin:0 1em}',
        '.__dexport_hint button[disabled]{opacity:.5}',
        '.__dexport_hint progress{margin-right:1em;vertical-align:middle;width:18em}',
    ]

    class Exporter {
        readonly timeSrv: TimeSrv
        oldRefreshRate?: string
        readonly elem: JQLite
        readonly loadProgress: JQLite
        readonly loadProgressText: JQLite
        readonly exportButton: JQLite
        readonly expandAllButton: JQLite
        loadProgressInterval?: number
        exportState: ExportState

        constructor() {
            const injector = $(document).injector();
            this.timeSrv = injector.get('timeSrv');
            const dashboard = this.timeSrv.dashboard;

            // construct the UI
            STYLES.forEach(style => FIRST_STYLESHEET.insertRule(style, 0));
            this.loadProgressText = $('<span>');
            this.loadProgress = $('<progress value="0">');
            this.exportButton = $('<button style="font-weight:bold">').text(L.export).one('click', () => this.startExport());
            this.expandAllButton = $('<button>').text(L.expandAllPanels).on('click', () => dashboard.expandRows());
            // note: order like [------- (Cancel) (Export)] on macOS,
            //              and [--- (Export) (Cancel) ---] everywhere else.
            const flexStyle = /^(Mac|iP)/.test(navigator.platform) ? 'flex-direction:row-reverse' : 'justify-content:center';
            this.elem = $('<div class="__dexport_hint">').append(
                $('<p>').text(L.exportTitle).append($('<strong>').text(dashboard.title)),
                $('<p>').text(L.exportSubtitle).append(this.expandAllButton),
                $('<p style="font-size:.8em">').append(this.loadProgress, this.loadProgressText),
                $(`<p style="display:flex;${flexStyle}">`).append(
                    this.exportButton,
                    $('<button>').text(L.cancel).on('click', () => this.closeUI()),
                ),
            );

            this.exportState = ExportState.idle;
            $(document.body).append(this.elem);

            // refresh the panel progress regularly
            this.loadProgressInterval = setInterval(() => this.updateProgress(), 500);
        }

        updateProgress(): void {
            const spinners = $('.panel-header');
            const total = spinners.length;

            // enable the export button if we are waiting too long; or immediately export if everything is loaded.
            if (this.exportState === ExportState.idle) {
                this.loadProgressText.text(`${L.numberOfPanels}: ${total}`)
                this.exportButton.prop('disabled', total === 0);
            } else {
                const ready = total - spinners.find('.panel-loading:visible').length;
                const percent = (100 * ready / total || 0).toFixed(1);
                this.loadProgressText.text(`${L.loadingPanels}: ${percent}% (${ready}/${total})`);
                this.loadProgress.prop({'value': ready, 'max': total});

                if (ready >= total) {
                    this.doExport();
                } else if (this.exportState !== ExportState.shortcut) {
                    this.exportState = ExportState.shortcut;
                    this.exportButton.prop('disabled', false).text(L.exportImmediately).one('click', () => this.doExport());
                }
            }
        }

        startExport(): void {
            const dashboard = this.timeSrv.dashboard;
            this.exportState = ExportState.started;
            this.exportButton.prop('disabled', true);

            // disable auto-refresh and ready the snapshot.
            this.oldRefreshRate = dashboard.refresh;
            this.timeSrv.setAutoRefresh();

            // remove the "ready" status of all spinners
            $('.panel-loading').removeClass('ng-hide');

            dashboard.snapshot = {timestamp: new Date()};
            dashboard.startRefresh();
        }

        doExport(): void {
            const dashboard = this.timeSrv.dashboard;

            // for the second click, take the dashboard snapshot.
            const timestamp = dashboard.snapshot!.timestamp.toISOString();
            const clone = dashboard.getSaveModelClone();
            const title = clone['title'];
            clone['time'] = this.timeSrv.timeRange();
            clone['id'] = null;
            clone['uid'] = null;
            clone['title'] = `${title} (exported at ${timestamp})`;

            const snapshot = {
                'meta': {
                    'isSnapshot': true,
                    'type': 'snapshot',
                    'expires': '9999-12-31T23:59:59Z',
                    'created': timestamp,
                },
                'dashboard': clone,
                'overwrite': true,
            };
            this.closeUI(); // cleanup immediately.

            // bring up the "save" dialog.
            const blob = new Blob([JSON.stringify(snapshot)], {type: 'application/json'});
            const url = URL.createObjectURL(blob);
            const a = $('<a>');
            a.prop({'href': url, 'download': `${title}_${timestamp}.json`});
            $(document.body).append(a);
            a[0].click();
            setTimeout(() => {
                URL.revokeObjectURL(url);
                a.remove();
            }, 0);

        }

        closeUI(): void {
            // scrub snapshot data
            const dashboard = this.timeSrv.dashboard;
            delete dashboard.snapshot;
            dashboard.forEachPanel(panel => delete panel.snapshotData);
            dashboard.annotations.list.forEach(anno => delete anno.snapshotData);

            // re-enable auto-refresh
            if (this.oldRefreshRate) {
                this.timeSrv.setAutoRefresh(this.oldRefreshRate);
            }
            clearInterval(this.loadProgressInterval);

            // remove the UI
            this.elem.remove();
            STYLES.forEach(() => FIRST_STYLESHEET.deleteRule(0));
        }
    }

    new Exporter();
} catch(e) {
    console.error(e);
    alert(L.unsupported);
}
