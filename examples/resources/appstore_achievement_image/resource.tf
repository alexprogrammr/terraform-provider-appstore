# Manage game center achievement localization image.
resource "appstore_achievement_image" "en-US" {
  achievement_localization_id = "<identifier of the achievement localization>"
  file                        = "img.png"
}
