package utils

import (
	"fmt"
	"reflect"
)

func CheckAndUpdateDomainAccordingToDTO(modifiableFields []string, entity interface{}, dto interface{}) (interface{}, error) {
	entityValue := reflect.ValueOf(entity).Elem()
	dtoValue := reflect.ValueOf(dto).Elem()

	for i := 0; i < dtoValue.NumField(); i++ {
		dtoField := dtoValue.Type().Field(i).Name
		dtoFieldValue := dtoValue.Field(i)

		entityField := entityValue.FieldByName(dtoField)

		// Eğer entity'de DTO'daki field yoksa, atla
		if !entityField.IsValid() {
			continue
		}

		// Eğer DTO'daki değer entity'deki değerle aynıysa atla
		if reflect.DeepEqual(entityField.Interface(), dtoFieldValue.Interface()) {
			continue
		}

		// Eğer DTO'daki alan entity'deki alandan farklıysa, modifiable listeye bak
		if isFieldModifiable(dtoField, modifiableFields) {
			// Alan modifiye edilebilirse entity'yi DTO'daki değerle güncelle
			if entityField.CanSet() {
				entityField.Set(dtoFieldValue)
			} else {
				return nil, fmt.Errorf("field %s cannot be updated", dtoField)
			}
		} else {
			// Eğer alan modifiye edilemiyorsa hata döndür
			return nil, fmt.Errorf("field %s is not modifiable", dtoField)
		}
	}

	// Her şey yolundaysa güncellenmiş entity'yi ve nil hatayı döndür
	return entity, nil
}

// Helper function to check if a field is modifiable
func isFieldModifiable(field string, modifiableFields []string) bool {
	for _, f := range modifiableFields {
		if f == field {
			return true
		}
	}
	return false
}
