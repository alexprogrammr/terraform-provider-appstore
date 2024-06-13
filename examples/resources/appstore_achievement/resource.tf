# Manage App Store Connect achievements.
resource "appstore_achievement" "test" {
  game_center_id     = "497799835"
  reference_name     = "Example Achievement"
  vendor_id          = "com.example.test"
  points             = 10
  repeatable         = false
  show_before_earned = false
}
