# rclip-client
Multiplatform client for RClip

## Installation
```
git clone https://github.com/NightWolf007/rclip-client.git
cd rclip-client
dep ensure
sudo go build -o /usr/local/bin/rclip-client
sudo mkdir /etc/rclip
sudo cp rclip-client.sample.yaml /etc/rclip/rclip-client.yaml
cp com.rclip.RClipClient.plist ~/Library/LaunchAgents/com.rclip.RClipClient.plist
launchctl load ~/Library/LaunchAgents/com.rclip.RClipClient.plist
```
