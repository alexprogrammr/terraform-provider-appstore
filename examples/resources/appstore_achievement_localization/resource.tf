# Manage game center achievement localization.
resource "appstore_achievement_localization" "en-US" {
  achievement_id            = "5ade5e98-7b45-42f9-a928-b513bf9fc279"
  locale                    = "en-US"
  name                      = "Test Achievement"
  before_earned_description = "Before earned description"
  after_earned_description  = "After earned description"
}
