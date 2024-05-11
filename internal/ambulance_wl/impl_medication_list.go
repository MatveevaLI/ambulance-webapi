package ambulance_wl

import (
	"net/http"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_waiting_list.go

// CreateMedicationListEntry - Saves new entry into waiting list
func (this *implMedicationListAPI) CreateMedicationListEntry(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry MedicationListEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(ambulance.MedicationList, func(waiting MedicationListEntry) bool {
			return entry.Id == waiting.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		ambulance.MedicationList = append(ambulance.MedicationList, entry)

		// entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(ambulance.MedicationList, func(waiting MedicationListEntry) bool {
			return entry.Id == waiting.Id
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return ambulance, ambulance.MedicationList[entryIndx], http.StatusOK
	})
}

// DeleteMedicationListEntry - Deletes specific entry
func (this *implMedicationListAPI) DeleteMedicationListEntry(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := ctx.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicationList, func(waiting MedicationListEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		ambulance.MedicationList = append(ambulance.MedicationList[:entryIndx], ambulance.MedicationList[entryIndx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}

// GetMedicationListEntries - Provides the ambulance medication list
func (this *implMedicationListAPI) GetMedicationListEntries(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		result := ambulance.MedicationList
		if result == nil {
			result = []MedicationListEntry{}
		}
		// return nil ambulance - no need to update it in db
		return nil, result, http.StatusOK
	})
}

// GetMedicationListEntry - Provides details about waiting list entry
func (this *implMedicationListAPI) GetMedicationListEntry(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := ctx.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicationList, func(waiting MedicationListEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		// return nil ambulance - no need to update it in db
		return nil, ambulance.MedicationList[entryIndx], http.StatusOK
	})
}

// UpdateMedicationListEntry - Updates specific entry
func (this *implMedicationListAPI) UpdateMedicationListEntry(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
        var entry MedicationListEntry

        if err := c.ShouldBindJSON(&entry); err != nil {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Invalid request body",
                "error":   err.Error(),
            }, http.StatusBadRequest
        }

        entryId := ctx.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(ambulance.MedicationList, func(waiting MedicationListEntry) bool {
            return entryId == waiting.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        if entry.Id != "" {
            ambulance.MedicationList[entryIndx].Id = entry.Id
        }

		if entry.Name != "" {
			ambulance.MedicationList[entryIndx].Name = entry.Name
		}

		if entry.Dosage != "" {
			ambulance.MedicationList[entryIndx].Dosage = entry.Dosage	
		}

		if entry.FrequencyPerDay != 0 {
			ambulance.MedicationList[entryIndx].FrequencyPerDay = entry.FrequencyPerDay
		}

		if entry.ExpirationDate != "" {
			ambulance.MedicationList[entryIndx].ExpirationDate = entry.ExpirationDate
		}

        return ambulance, ambulance.MedicationList[entryIndx], http.StatusOK
    })
}
