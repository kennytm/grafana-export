# Manual snapshot exporting

This method is suitable when you cannot use the browser's development tools.

1. Login to Grafana

2. Open the dashboard you wish export

3. Click "Share dashboard" on the toolbar on top-right corner

    ![The "Share dashbaord" button](../docs/manual-1.png)

4. Choose the "Snapshot" tab, and click "Local Snapshot"

    ![The "Local snapshot" button](../docs/manual-4.png)

5. Wait until the link is available

    ![Result of clicking "Local Snapshot"](../docs/manual-2.png)

6. Copy the link, and then replace the `dashboard/snapshot` part in the URL by `api/snapshots`, and visit the new link.

    ![localhost:3000/dashboard/snapshot/...](../docs/manual-3a.png)<br>
    ↓<br>
    ![localhost:3000/api/snapshots/...](../docs/manual-3b.png)

7. Save the current web page (<kbd>Ctrl</kbd>+<kbd>S</kbd>) as a JSON file.

8. Repeat from step 2 for other dashboards.

