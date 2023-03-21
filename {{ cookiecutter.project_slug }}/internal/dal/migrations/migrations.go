package migrations

import (
	"cp/internal/dal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		model.Customer{},
		model.CustomerCalculation{}, // NEW
		model.CustomerSalesStageRecord{},

		model.Conversation{},
		model.ConversationCalculation{}, // NEW

		model.BIDConfig{},
		model.StaffConfig{},

		// 重构标签树
		model.BlackList{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}
	// m := gormigrate.New(db, &migrationOpts, migrations)
	// if err := m.Migrate(); err != nil {
	// 	log.Fatal().Err(err).Msg("could not migrate")
	// }
}

// var migrations = []*gormigrate.Migration{
// 	// :迁移conversation表的一次性操作
// 	{
// 		ID: "2022052501",
// 		Migrate: func(tx *gorm.DB) error {
// 			conversation := query.Use(tx).Conversation
// 			if _, err := conversation.WithContext(context.Background()).
// 				Where(conversation.CallID.IsNull()).
// 				Update(conversation.CallID, conversation.ID); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	},
// 	{
// 		ID: "2022081201",
// 		Migrate: func(tx *gorm.DB) error {
// 			migrator := tx.Migrator()
// 			if migrator.HasColumn(model.Customer{}, "properties") {
// 				if err := tx.Exec("ALTER TABLE `customers` MODIFY COLUMN `properties` JSON COMMENT 'deprecated'").Error; err != nil {
// 					return err
// 				}
// 			}
// 			if migrator.HasColumn(model.Customer{}, "used_names") {
// 				if err := tx.Exec("ALTER TABLE `customers` MODIFY COLUMN `used_names` JSON COMMENT 'deprecated'").Error; err != nil {
// 					return err
// 				}
// 			}
// 			if migrator.HasColumn(model.Customer{}, "phone_suffix") {
// 				if err := tx.Exec("ALTER TABLE `customers` MODIFY COLUMN `phone_suffix` char(4) COMMENT 'deprecated'").Error; err != nil {
// 					return err
// 				}
// 			}
// 			return nil
// 		},
// 	},
// }
