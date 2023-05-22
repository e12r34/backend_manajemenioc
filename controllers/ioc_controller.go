package controllers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fiberioc/configs"
	"fiberioc/libs"
	"fiberioc/models"
	"fiberioc/responses"

	"log"
	"net/http"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var iocCollection *mongo.Collection = configs.GetCollection(configs.DB, "cl_ioc")
var pulseCollection *mongo.Collection = configs.GetCollection(configs.DB, "cl_pulse")

// CreateIoc Ioc
// @Summary      Create New Ioc
// @Description  Create New Ioc
// @Tags         iocs
// @Accept       json
// @Param request body models.Ioc_input true "query params"
// @Produce      json
// @Success      200  {object}  models.Ioc
// @Router       /ioc [post]
func CreateIoc(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var ioc models.Ioc_input
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&ioc); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&ioc); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}

	result, err := iocCollection.InsertOne(ctx, ioc)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: result})
}

// GetAIoc 1 Ioc
// @Summary      Get 1 Ioc
// @Description  Get 1 Ioc By id
// @Tags         iocs
// @Accept       json
// @Param        id   path      string  true  "id ioc"
// @Produce      json
// @Success      200  {object}  models.Ioc
// @Router       /ioc/{id} [get]
func GetAIoc(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var ioc models.Ioc
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	err := iocCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&ioc)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: ioc})
}

// Edit Ioc
// @Summary      Edit Ioc
// @Description  Edit Existing Ioc
// @Tags         iocs
// @Accept       json
// @Param request body models.Ioc_input true "query params"
// @Produce      json
// @Success      200  {object}  models.Ioc
// @Router       /ioc/{id} [put]
func EditAIoc(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var ioc_edit models.Ioc_Edit
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	//validate the request body
	if err := c.BodyParser(&ioc_edit); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error1", Data: err.Error()})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&ioc_edit); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error2", Data: validationErr.Error()})
	}

	var myMap map[string]interface{}
	data, _ := json.Marshal(ioc_edit)
	json.Unmarshal(data, &myMap)
	a := bson.M(myMap)
	result, err := iocCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": a})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error3", Data: err.Error()})
	}

	//get updated user details
	var updatedIoc models.Ioc
	if result.MatchedCount == 1 {
		err := iocCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedIoc)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error4", Data: err.Error()})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: updatedIoc})
}

func DeleteAIoc(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: "User with specified ID not found!"},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

// GetAllIoc Iocs
// @Summary      Show an account
// @Description  get all data
// @Tags         iocs
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Ioc
// @Router       /iocs [get]
func GetAllIocs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var iocs []models.Ioc
	defer cancel()

	results, err := iocCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var ioc models.Ioc
		if err = results.Decode(&ioc); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		iocs = append(iocs, ioc)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: iocs},
	)
}

// GetAllIoc Iocs
// @Summary      Show data
// @Description  get data as csv
// @Tags         iocs
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Ioc
// @Router       /iocs/csv [get]
func GetAllIocsCSV(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var iocs []models.Ioc
	defer cancel()

	results, err := iocCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var ioc models.Ioc
		if err = results.Decode(&ioc); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		iocs = append(iocs, ioc)
	}
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="test.csv"`)
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return gocsv.Marshal(iocs, c)
}

// GetAllIoc Iocs
// @Summary      Show some
// @Description  get some data
// @Tags         iocs
// @Accept       json
// @Produce      json
// @Param request body models.Ioc_Get true "query params"
// @Success      200  {object}  models.Ioc
// @Router       /iocs/some [get]
func GetSomeIoc(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var idioc models.Ioc_Get
	var iocs []models.Ioc
	defer cancel()

	if err := c.BodyParser(&idioc); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	filter := bson.M{"_id": bson.M{"$in": idioc.Id}}
	results, err := iocCollection.Find(ctx, filter)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}
	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var ioc models.Ioc
		if err = results.Decode(&ioc); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		iocs = append(iocs, ioc)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: iocs},
	)
}

func GetSomeIocCSV(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var idioc models.Ioc_Get
	var iocs []models.Ioc
	defer cancel()

	if err := c.BodyParser(&idioc); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error1", Data: &fiber.Map{"data": err.Error()}})
	}

	filter := bson.M{"_id": bson.M{"$in": idioc.Id}}

	results, err := iocCollection.Find(ctx, filter)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}
	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var ioc models.Ioc
		if err = results.Decode(&ioc); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		iocs = append(iocs, ioc)
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="export_`+time.Now().Format("02-Jan-2006")+`.csv"`)
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return gocsv.Marshal(iocs, c)
}

func DuaCSV(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var ipZIPFile bytes.Buffer
	// Create a new CSV writer for the "people" collection
	ipCSVWriter := zip.NewWriter(&ipZIPFile)

	// Write the header row for the "people" CSV file
	cursor, err := iocCollection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"indicator": 1, "otx_id": 1, "source": 1, "description": 1, "type": 1, "_id": 1}))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var csvDataForCollection bytes.Buffer
	writer := csv.NewWriter(&csvDataForCollection)
	writer.Comma = ';'
	// Write header row

	err = writer.Write([]string{"Id", "Type", "Indicator", "Description", "Source", "Source_Hash"})
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		var ipField models.Ioc
		err = cursor.Decode(&ipField)

		if err != nil {
			log.Fatal(err)
		}

		err = writer.Write([]string{ipField.Id.Hex(), ipField.Type, ipField.Indicator, ipField.Description, ipField.Source, ipField.OTX_Hash})
		if err != nil {
			log.Fatal(err)
		}

	}
	writer.Flush()

	// Add CSV file to zip
	zipFile, err := ipCSVWriter.Create("ip.csv")
	if err != nil {
		return err
	}
	_, err = zipFile.Write(csvDataForCollection.Bytes())
	if err != nil {
		return err
	}

	// Close zip writer
	err = ipCSVWriter.Close()
	if err != nil {
		return err
	}

	// Set response headers
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=export.zip")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	// Send zip file as response
	_, err = c.Write(ipZIPFile.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func ExportThorAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var ipZIPFile bytes.Buffer

	ipZIPWriter := zip.NewWriter(&ipZIPFile)

	filter := bson.M{
		"type": bson.M{
			"$regex": "^FileHash*",
		},
	}
	_, err := libs.ThorMakeHash(ipZIPWriter, ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	filter = bson.M{
		"$or": bson.A{
			bson.M{"type": "IPv4"},
			bson.M{"type": "hostname"},
			bson.M{"type": "domain"},
			bson.M{"type": "URL"},
		},
	}
	_, err = libs.ThorMakeC2C(ipZIPWriter, ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	filter = bson.M{"type": "YARA"}
	_, err = libs.ThorMakeYara(ipZIPWriter, ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Close zip writer
	err = ipZIPWriter.Close()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=export-thor.zip")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	// Send zip file as response
	_, err = c.Write(ipZIPFile.Bytes())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return nil
}

func ExportThorSome(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	var idioc models.Ioc_Get

	if err := c.BodyParser(&idioc); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var ipZIPFile bytes.Buffer

	ipZIPWriter := zip.NewWriter(&ipZIPFile)

	//filter untuk include apakah id ada di db
	filter1 := bson.M{
		"_id": bson.M{
			"$in": idioc.Id,
		},
	}

	//filter untuk regex mengambil tipe hash
	filter2 := bson.M{
		"type": bson.M{
			"$regex": "^FileHash*",
		},
	}

	//define filter all
	filter := bson.M{
		"$and": bson.A{
			filter1,
			filter2,
		},
	}

	_, err := libs.ThorMakeHash(ipZIPWriter, ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	//update filter regex hash, menjadi ioc
	filter2 = bson.M{
		"$or": bson.A{
			bson.M{"type": "IPv4"},
			bson.M{"type": "hostname"},
			bson.M{"type": "domain"},
			bson.M{"type": "URL"},
		},
	}

	// update filter al
	filter = bson.M{
		"$and": bson.A{
			filter1,
			filter2,
		},
	}
	_, err = libs.ThorMakeC2C(ipZIPWriter, ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	filter2 = bson.M{"type": "YARA"}

	// update filter al
	filter = bson.M{
		"$and": bson.A{
			filter1,
			filter2,
		},
	}
	_, err = libs.ThorMakeYara(ipZIPWriter, ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Close zip writer
	err = ipZIPWriter.Close()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=export-thor.zip")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	// Send zip file as response
	_, err = c.Write(ipZIPFile.Bytes())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return nil
	// return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
}

func ExportMgAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": bson.A{
			bson.M{"type": "IPv4"},
			bson.M{"type": "hostname"},
			bson.M{"type": "domain"},
			bson.M{"type": "URL"},
		},
	}

	JsonMGRules, err := libs.MgMakeRules(ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", "attachment; filename=export-MG.json")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return c.Status(fiber.StatusOK).Type("json").Send(JsonMGRules)
}

func ExportMgSome(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var idioc models.Ioc_Get

	if err := c.BodyParser(&idioc); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	filter1 := bson.M{
		"_id": bson.M{
			"$in": idioc.Id,
		},
	}

	//filter untuk regex mengambil tipe hash
	filter2 := bson.M{
		"$or": bson.A{
			bson.M{"type": "IPv4"},
			bson.M{"type": "hostname"},
			bson.M{"type": "domain"},
			bson.M{"type": "URL"},
		},
	}

	//define filter all
	filter := bson.M{
		"$and": bson.A{
			filter1,
			filter2,
		},
	}

	JsonMGRules, err := libs.MgMakeRules(ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", "attachment; filename=export-MG.json")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Status(fiber.StatusOK).Type("json").Send(JsonMGRules)
}

func ExportBAEAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": bson.A{
			bson.M{"type": "IPv4"},
			bson.M{"type": "hostname"},
			bson.M{"type": "domain"},
			bson.M{"type": "URL"},
		},
	}

	BAERules, err := libs.BAEMakeRules(ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/csv")
	c.Set("Content-Disposition", "attachment; filename=export-BAE.csv")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	// return c.Status(fiber.StatusOK).Type("csv").Send(BAERules)
	return c.SendStream(BAERules, int(BAERules.Len()))
}

func ExportBAESome(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var idioc models.Ioc_Get

	if err := c.BodyParser(&idioc); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	filter1 := bson.M{
		"_id": bson.M{
			"$in": idioc.Id,
		},
	}

	filter2 := bson.M{
		"$or": bson.A{
			bson.M{"type": "IPv4"},
			bson.M{"type": "hostname"},
			bson.M{"type": "domain"},
			bson.M{"type": "URL"},
		},
	}

	//define filter all
	filter := bson.M{
		"$and": bson.A{
			filter1,
			filter2,
		},
	}

	BAERules, err := libs.BAEMakeRules(ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/csv")
	c.Set("Content-Disposition", "attachment; filename=export-BAE.csv")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	// return c.Status(fiber.StatusOK).Type("csv").Send(BAERules)
	return c.SendStream(BAERules, int(BAERules.Len()))
}

func ExportSuricataAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"type": "IPv4",
	}

	SuricataRules, err := libs.SuricataMakeIPRep(ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=export-Suricata.zip")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	// return c.Status(fiber.StatusOK).Type("csv").Send(BAERules)
	_, err = c.Write(SuricataRules.Bytes())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return nil
}

func ExportSuricataSome(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var idioc models.Ioc_Get

	if err := c.BodyParser(&idioc); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	filter1 := bson.M{
		"_id": bson.M{
			"$in": idioc.Id,
		},
	}

	filter2 := bson.M{
		"type": "IPv4",
	}

	//define filter all
	filter := bson.M{
		"$and": bson.A{
			filter1,
			filter2,
		},
	}

	SuricataRules, err := libs.SuricataMakeIPRep(ctx, iocCollection, pulseCollection, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// Set response headers
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=export-Suricata.zip")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	// return c.Status(fiber.StatusOK).Type("csv").Send(BAERules)
	_, err = c.Write(SuricataRules.Bytes())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return nil
}
