package language

import (
	"errors"
	"fyssl/connection/actions/targets/reply/language/asterisk"
	"fyssl/connection/actions/targets/reply/language/types"
)

var keyToProcess = map[string]func (fullExpression string, buff *[]byte) (*[]byte, error) {
	"*": asterisk.ProcessAsterisk,
}

func Process(buff *[]byte, replyExpression string) (*[]byte, error) {
	magicSlice := makeMagicSlice(&replyExpression)
	var processedBuff []byte
	for len(*magicSlice) >= 2 {
		a, b, err := findInner(magicSlice)
		if err != nil {
			return buff, err
		}
		result, err := keyToProcess[a.Key](replyExpression[a.Location+1:b.Location], buff)
		if err != nil {
			return buff, err
		}
		processedBuff = append(processedBuff, *result...)
		magicSlice = removeMagic(magicSlice, a)
		magicSlice = removeMagic(magicSlice, b)
	}
	return &processedBuff, nil
}

func makeMagicSlice(replyExpression *string) *[]types.KeyToLocation {
	var magics []types.KeyToLocation
	for index, key := range *replyExpression {
		for magic, _ := range keyToProcess {
			if magic == string(key) {
				magics = append(magics, types.KeyToLocation{
					Key:      string(key),
					Location: index,
				})
			}
		}
	}
	return &magics
}


func findInner(magicSlice *[]types.KeyToLocation) (types.KeyToLocation, types.KeyToLocation, error) {
	for index, _ := range *magicSlice {
		if len(*magicSlice) != index - 1 && (*magicSlice)[index].Key == (*magicSlice)[index+1].Key {
			return (*magicSlice)[index], (*magicSlice)[index+1], nil
		}
	}
	return types.KeyToLocation{}, types.KeyToLocation{}, errors.New("no inner keys found")
}

func removeMagic(magicSlice *[]types.KeyToLocation, magic types.KeyToLocation) *[]types.KeyToLocation {
	var newSlice []types.KeyToLocation
	for _, element := range *magicSlice {
		if element != magic {
			newSlice = append(newSlice, element)
		}
	}
	return &newSlice
}