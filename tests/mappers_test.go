package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jental/webbugs-common-go/contracts"
	"github.com/jental/webbugs-common-go/mappers"
	"github.com/jental/webbugs-common-go/models"
)

func TestMapToData(t *testing.T) {
	jsonStr, err := ioutil.ReadFile("../resources/data0.test.resources.json")
	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
	}

	var contract contracts.DataContract
	err = json.Unmarshal(jsonStr, &contract)
	if err != nil {
		t.Log("error (2) should be nil", err)
		t.Fail()
	}

	field, components, err := mappers.MapToData(&contract)
	if err != nil {
		t.Log("error (3) should be nil", err)
		t.Fail()
	}

	log.Println(field, components)

	cell := field.Get(models.NewCoordinates(3, -8, 5))
	if cell == nil {
		t.Log("Cell [3,-8,5] is expected to be present")
		t.Fail()
	}

	if cell.CellType != models.CellTypeWall {
		t.Log("Cell [3,-8,5] has unexpected cell type", cell.CellType)
		t.Fail()
	}

	expectedUUID, err := uuid.Parse("a1322c04-1029-4bb4-92d0-34d865c90e46")
	if cell.PlayerID != expectedUUID {
		t.Log("Cell [3,-8,5] has unexpected player id", cell.PlayerID)
		t.Fail()
	}

	if cell.Crd != models.NewCoordinates(3, -8, 5) {
		t.Log("Cell [3,-8,5] has unexpected coordinates", cell.Crd)
		t.Fail()
	}

	if cell.Component == nil {
		t.Log("Cell [3,-8,5] is expected to have a component")
		t.FailNow()
	}
	componentID := cell.Component.ID

	component, succ := components.Get(componentID)
	if !succ {
		t.Log(fmt.Sprintf("Component [%d] is expected to be present", componentID))
		t.Fail()
	}

	if component.ID != componentID {
		t.Log(fmt.Sprintf("Component [%d] is expected to have an id 1", componentID))
		t.Fail()
	}

	if component.IsActive {
		t.Log(fmt.Sprintf("Component [%d] is expected to be inactive", componentID))
		t.Fail()
	}

	if len(component.Walls) != 13 {
		t.Log(fmt.Sprintf("Component [%d] is expected to have 13 walls", componentID))
		t.Fail()
	}
}

func TestMapFromData(t *testing.T) {
	field := models.NewField(3)
	crd := models.NewCoordinates(1, 2, 1)
	playerID, err := uuid.NewUUID()
	if err != nil {
		t.Log("Failed to generate layer id.")
		t.Fail()
	}
	cell := models.NewBugCell(playerID, crd, true)
	field.Set(crd, &cell)
	crd2 := models.NewCoordinates(2, 1, 2)
	walls := make([]*models.Cell, 1)
	component := models.NewComponent(true, walls)
	wall := models.NewWallCell(playerID, crd2, &component)
	walls[0] = &wall
	var components models.Components
	components.Set(&component)

	dataContract := mappers.MapFromData(&field, &components)
	t.Log(fmt.Sprintf("Contract [%v]", dataContract))

	dataContract.Field.Grid[0].Grid[169].PlayerID = uuid.MustParse("16639eed-df40-11eb-94e9-02004c4f4f50")

	j, err := json.Marshal(dataContract)
	if err != nil {
		t.Log(fmt.Sprintf("Error during serialization: %v", err))
		t.Fail()
	}

	expJson, err := ioutil.ReadFile("../resources/data1.test.resources.json")
	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
	}

	for i, ch := range j {
		if expJson[i] != ch {
			t.Log("different symbols:", i, expJson[i], ch)
			t.Log(string(j))
			t.Log(string(expJson))
			t.FailNow()
		}
	}
}
