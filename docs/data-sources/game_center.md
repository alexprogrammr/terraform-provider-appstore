---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "appstore_game_center Data Source - appstore"
subcategory: ""
description: |-
  Fetches Game Center information from the App Store Connect.
---

# appstore_game_center (Data Source)

Fetches Game Center information from the App Store Connect.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_id` (String) Identifier of the app to fetch Game Center information for.

### Read-Only

- `arcade_enabled` (Boolean) Indicates whether Game Center is enabled for the app on Apple Arcade.
- `challenge_enabled` (Boolean) Indicates whether Game Center challenges are enabled for the app.
- `id` (String) Identifier of the game center.