# TV Source Switch

Automatically switch to TV source when PC is turned on.

Uses the Samsung SmartThings API

## Installation

1. Copy compiled binary to `shell:startup` folder. 
2. Create a SmartThings personal access token at https://my.smartthings.com/tokens
3. Use the following curl request to determine your device ID:
   ```
   curl --location --request GET 'https://api.smartthings.com/v1/devices' \
   --header 'Authorization: Bearer <your token>'
   ```
4. Create a config file in `%HOME%\.config\tvsourceswitch\config.yaml` folder with the following contents
    ```yaml
    source: HDMI4
    smartthings_token: <your token>
    smartthings_device_id: <your device id>
    ```