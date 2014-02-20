package helpers

import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/database/models"

func GetCharactersForClient(c interfaces.Client) ([]*models.Character, error) {
	characters := make([]*models.Character, 0)

	charactersRaw, err := c.Server().Database().Select(&models.Character{}, "SELECT * FROM characters WHERE account_id = ?", c.Account().AccountID)
	if err != nil {
		return nil, err
	}

	for _, characterRaw := range charactersRaw {
		character, ok := characterRaw.(*models.Character)
		if ok {
			characters = append(characters, character)
		}
	}

	return characters, nil
}

func IsCharacterNameTaken(c interfaces.Client, characterName string) (bool, error) {
	characterCount, err := c.Server().Database().SelectInt("SELECT COUNT(*) FROM characters WHERE first_name = ?", characterName)
	if err != nil {
		return false, err
	}

	return characterCount > 0, nil
}
