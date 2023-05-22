package libs

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fiberioc/models"
	"io"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ThorMakeHash(ipZIPWriter *zip.Writer, ctx context.Context, Collection *mongo.Collection, pulseCollection *mongo.Collection, filter primitive.M) (io.Writer, error) {
	cursor, err := Collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"indicator": 1, "otx_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Buffer bytes.Buffer

	Buffer.WriteString("#\n# THOR CUSTOM HASH IOCs\n# This file contains MD5, SHA1 and SHA256 hashes and a short info like file name\n# or hash origin\n#\n# Important: Rename this template file from .txt.template to .txt \n#\n# FORMAT -----------------------------------------------------------------------\n#\n# MD5;COMMENT\n# SHA1;COMMENT\n# SHA256;COMMENT\n#\n# EXAMPLES ---------------------------------------------------------------------\n#\n# 0c2674c3a97c53082187d930efb645c2;DEEP PANDA Sakula Malware - http://goo.gl/R3e6eG\n# 000c907d39924de62b5891f8d0e03116;The Darkhotel APT http://goo.gl/DuS7WS\n# c03318cb12b827c03d556c8747b1e323225df97bdc4258c2756b0d6a4fd52b47;Operation SMN Hashes http://goo.gl/bfmF8B - Zxshell\n")

	mapName := make(map[string]string)
	for cursor.Next(context.Background()) {
		var Ioc models.Ioc_Thor
		err := cursor.Decode(&Ioc)
		if err != nil {
			return nil, err
		}
		value, isNotNew := mapName[Ioc.Id]

		if isNotNew {
			Ioc.Name = value
		} else {
			filterPulse := bson.M{
				"id": Ioc.Id,
			}
			var result bson.M
			err := pulseCollection.FindOne(context.Background(), filterPulse, options.FindOne().SetProjection(bson.M{"name": 1})).Decode(&result)
			if err != nil {
				return nil, err
			}

			Ioc.Name = result["name"].(string)
			mapName[Ioc.Id] = Ioc.Name
		}
		Buffer.WriteString(Ioc.Indicator + ";" + Ioc.Name + "\n")

	}

	zipFile, err := ipZIPWriter.Create("custom-hash-iocs.txt")
	if err != nil {
		return nil, err
	}

	_, err = zipFile.Write(Buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func ThorMakeC2C(ipZIPWriter *zip.Writer, ctx context.Context, Collection *mongo.Collection, pulseCollection *mongo.Collection, filter primitive.M) (io.Writer, error) {

	cursor, err := Collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"indicator": 1, "otx_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Buffer bytes.Buffer

	Buffer.WriteString("#\n# THOR CUSTOM C2 IOCs\n#\n# Important: Activate this template by renaming it from .txt.template to .txt \n#\n# FORMAT -----------------------------------------------------------------------\n#\n# # Comment\n# IP\n# FQDN\n#\n# EXAMPLES ---------------------------------------------------------------------\n#\n## Case 44 C2 Server\n#mastermind.eu\n#googleaccountservices.com\n#89.22.123.12\n")

	for cursor.Next(context.Background()) {
		var Ioc models.Ioc_Thor
		err := cursor.Decode(&Ioc)
		if err != nil {
			return nil, err
		}
		Buffer.WriteString(Ioc.Indicator + "\n")

	}
	zipFile, err := ipZIPWriter.Create("custom-c2-domains.txt")
	if err != nil {
		return nil, err
	}

	_, err = zipFile.Write(Buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func ThorMakeYara(ipZIPWriter *zip.Writer, ctx context.Context, Collection *mongo.Collection, pulseCollection *mongo.Collection, filter primitive.M) (io.Writer, error) {

	cursor, err := Collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"content": 1, "otx_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Buffer bytes.Buffer

	// Buffer.WriteString("#\n# THOR CUSTOM C2 IOCs\n#\n# Important: Activate this template by renaming it from .txt.template to .txt \n#\n# FORMAT -----------------------------------------------------------------------\n#\n# # Comment\n# IP\n# FQDN\n#\n# EXAMPLES ---------------------------------------------------------------------\n#\n## Case 44 C2 Server\n#mastermind.eu\n#googleaccountservices.com\n#89.22.123.12\n")

	for cursor.Next(context.Background()) {
		var Ioc models.Ioc_Thor
		err := cursor.Decode(&Ioc)
		if err != nil {
			return nil, err
		}
		Buffer.WriteString(Ioc.Content + "\n\n")

	}
	zipFile, err := ipZIPWriter.Create("custom-yara.yar")
	if err != nil {
		return nil, err
	}

	_, err = zipFile.Write(Buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func MgMakeRules(ctx context.Context, Collection *mongo.Collection, pulseCollection *mongo.Collection, filter primitive.M) ([]byte, error) {
	option := options.Find().SetProjection(bson.M{"indicator": 1, "otx_id": 1, "type": 1}).SetSort(bson.D{{"otx_id", 1}})
	cursor, err := Collection.Find(ctx, filter, option)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mapPulse := make(map[string]string)
	result := make(map[string]any)
	urutan := make(map[string]int)
	i := 0
	for cursor.Next(context.Background()) {
		Ioc := make(map[string]interface{})
		var fromDB models.Ioc_MG_BAE
		err := cursor.Decode(&fromDB)

		if err != nil {
			return nil, err
		}
		value, isNotNew := mapPulse[fromDB.Id]

		if isNotNew {
			urutan[fromDB.Id]++
			Ioc["malware_name"] = value + "-" + strconv.FormatInt(int64(urutan[fromDB.Id]), 10)

		} else {
			filterPulse := bson.M{
				"id": fromDB.Id,
			}
			var result bson.M
			option2 := options.FindOne().SetProjection(bson.M{"name": 1})
			err := pulseCollection.FindOne(context.Background(), filterPulse, option2).Decode(&result)
			if err != nil {
				return nil, err
			}
			nama := result["name"].(string)
			if len(nama) > 60 {
				nama = nama[:60]
			}
			urutan[fromDB.Id] = 1
			mapPulse[fromDB.Id] = nama + "-" + strconv.FormatInt(int64(urutan[fromDB.Id]), 10)
			Ioc["malware_name"] = mapPulse[fromDB.Id]

		}

		Ioc["malware_name"] = strings.ReplaceAll(Ioc["malware_name"].(string), "`", "")
		Ioc["is_tpd"] = "f"
		Ioc["desc"] = ""
		Ioc["category"] = "Malware-Other Malware"
		Ioc["active"] = "t"
		if fromDB.Type == "hostname" || fromDB.Type == "domain" {
			Ioc["host"] = fromDB.Indicator
			Ioc["type"] = "host"
		} else if fromDB.Type == "IPv4" {
			Ioc["dip"] = fromDB.Indicator
			Ioc["type"] = "dip"
		} else if fromDB.Type == "URL" {
			URL := strings.SplitN(strings.ReplaceAll(strings.ReplaceAll(fromDB.Indicator, "https://", ""), "http://", ""), "/", 2)
			Ioc["host"] = URL[0]
			Ioc["uri"] = URL[1]
			Ioc["type"] = "host:dport:uri"

			if strings.Contains(fromDB.Indicator, "https://") {
				Ioc["dport"] = 443
			} else {
				Ioc["dport"] = 80
			}
		}

		result[strconv.FormatInt(int64(len(result)), 10)] = Ioc
		i++

	}

	// Marshal the map as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func BAEMakeRules(ctx context.Context, Collection *mongo.Collection, pulseCollection *mongo.Collection, filter primitive.M) (*bytes.Buffer, error) {
	option := options.Find().SetProjection(bson.M{"indicator": 1, "otx_id": 1, "type": 1}).SetSort(bson.D{{"otx_id", 1}})
	cursor, err := Collection.Find(ctx, filter, option)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// header

	header := "id,trafficCaptured,remarks,threatActor,source,createdDate,reportNumber,threatType,additionDate,lastUpdatedDate,ruleValue,wildcardPrefix,wildcardSuffix,priority,confidence,matchThreshold,type,case,campaign,taskStatus,status.name"

	templates := ",false,,$_APT_$,User,,,APT,,,$_IoC_$,true,true,$_Priority_$,$_Confidence_$,,$_Category_$,,,Active,"
	format_apt := "$_APT_$"
	format_ioc := "$_IoC_$"
	format_priority := "$_Priority_$"
	format_confidence := "$_Confidence_$"
	format_category := "$_Category_$"
	new_line := "\n"
	var Buffer bytes.Buffer
	Buffer.WriteString(header + new_line)
	mapPulse := make(map[string]string)
	mapTLP := make(map[string]string)
	for cursor.Next(context.Background()) {
		template := templates
		var fromDB models.Ioc_MG_BAE
		err := cursor.Decode(&fromDB)
		template = strings.ReplaceAll(template, format_ioc, fromDB.Indicator)

		if fromDB.Type == "hostname" || fromDB.Type == "domain" || fromDB.Type == "URL" {
			template = strings.ReplaceAll(template, format_category, "FQDN")
		} else if fromDB.Type == "IPv4" {
			template = strings.ReplaceAll(template, format_category, fromDB.Type)
		}
		if err != nil {
			return nil, err
		}
		value, isNotNew := mapPulse[fromDB.Id]

		if isNotNew {
			template = strings.ReplaceAll(template, format_apt, value)

		} else {
			filterPulse := bson.M{
				"id": fromDB.Id,
			}
			var result bson.M
			option2 := options.FindOne().SetProjection(bson.M{"name": 1, "tlp": 1})
			err := pulseCollection.FindOne(context.Background(), filterPulse, option2).Decode(&result)
			if err != nil {
				return nil, err
			}
			mapTLP[fromDB.Id] = result["tlp"].(string)
			mapPulse[fromDB.Id] = strings.ReplaceAll(result["name"].(string), ",", "-")
			value = mapPulse[fromDB.Id]
			template = strings.ReplaceAll(template, format_apt, value)
		}
		template = strings.ReplaceAll(template, format_priority, "5")
		template = strings.ReplaceAll(template, format_confidence, "5")
		// if mapTLP[fromDB.Id] == "white" {
		// 	template = strings.ReplaceAll(template, format_priority, "1")
		// 	template = strings.ReplaceAll(template, format_confidence, "1")
		// } else if mapTLP[fromDB.Id] == "green" {
		// 	template = strings.ReplaceAll(template, format_priority, "2")
		// 	template = strings.ReplaceAll(template, format_confidence, "2")
		// } else if mapTLP[fromDB.Id] == "amber" {
		// 	template = strings.ReplaceAll(template, format_priority, "3")
		// 	template = strings.ReplaceAll(template, format_confidence, "3")
		// } else if mapTLP[fromDB.Id] == "amber+strict" {
		// 	template = strings.ReplaceAll(template, format_priority, "4")
		// 	template = strings.ReplaceAll(template, format_confidence, "4")
		// } else if mapTLP[fromDB.Id] == "red" {
		// 	template = strings.ReplaceAll(template, format_priority, "5")
		// 	template = strings.ReplaceAll(template, format_confidence, "5")
		// }
		Buffer.WriteString(template + new_line)
	}

	return &Buffer, nil
}

func SuricataMakeIPRep(ctx context.Context, Collection *mongo.Collection, pulseCollection *mongo.Collection, filter primitive.M) (*bytes.Buffer, error) {
	var ipZIPFile bytes.Buffer

	ipZIPWriter := zip.NewWriter(&ipZIPFile)

	cursor, err := Collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"indicator": 1, "otx_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	sid := 99990000
	TemplateRules := "alert ip any any -> any any (msg:\"$_nama_$\"; flow:stateless; iprep:any,$_pulse_$,>,30; sid:$_sid_$; reference:$_ref_$; rev:1;)"
	TemplateReference := "$_ref_$"
	TemplateNama := "$_nama_$"
	TemplatePulse := "$_pulse_$"
	TemplateSid := "$_sid_$"
	new_line := "\n"
	var BufferIPRepRules bytes.Buffer
	var BufferRepList bytes.Buffer
	var BufferCategory bytes.Buffer
	urutan := 1
	idCategory := make(map[string]string)
	mapName := make(map[string]string)
	for cursor.Next(context.Background()) {
		var Ioc models.Ioc_Thor
		err := cursor.Decode(&Ioc)
		if err != nil {
			return nil, err
		}
		value, isNotNew := mapName[Ioc.Id]

		if isNotNew {
			Ioc.Name = value
		} else {
			filterPulse := bson.M{
				"id": Ioc.Id,
			}
			var result bson.M
			err := pulseCollection.FindOne(context.Background(), filterPulse, options.FindOne().SetProjection(bson.M{"name": 1})).Decode(&result)
			if err != nil {
				return nil, err
			}

			Ioc.Name = result["name"].(string)
			mapName[Ioc.Id] = Ioc.Name
			idCategory[Ioc.Id] = strconv.FormatInt(int64(urutan), 10)
			BufferCategory.WriteString(strconv.FormatInt(int64(urutan), 10) + "," + Ioc.Id + "," + cleansing(Ioc.Name) + "\n")

			Template := strings.ReplaceAll(TemplateRules, TemplateNama, "Bad Reputation IP Connection from Pulse OTX "+Ioc.Id)
			Template = strings.ReplaceAll(Template, TemplateSid, strconv.FormatInt(int64(sid), 10))
			sid++
			Template = strings.ReplaceAll(Template, TemplatePulse, Ioc.Id)
			Template = strings.ReplaceAll(Template, TemplateReference, "url,https://otx.alienvault.com/pulse/"+Ioc.Id)
			BufferIPRepRules.WriteString(Template + new_line)

			urutan++
		}
		BufferRepList.WriteString(Ioc.Indicator + "," + idCategory[Ioc.Id] + ",127" + new_line)
		// Buffer.WriteString(Ioc.Indicator + ";" + Ioc.Name + "\n")

	}

	FileCategory, err := ipZIPWriter.Create("category-otx.txt")
	if err != nil {
		return nil, err
	}

	_, err = FileCategory.Write(BufferCategory.Bytes())
	if err != nil {
		return nil, err
	}

	FileIPRepRules, err := ipZIPWriter.Create("iprep-otx.rules")
	if err != nil {
		return nil, err
	}

	_, err = FileIPRepRules.Write(BufferIPRepRules.Bytes())
	if err != nil {
		return nil, err
	}

	FileRepList, err := ipZIPWriter.Create("iprep-otx.list")
	if err != nil {
		return nil, err
	}

	_, err = FileRepList.Write(BufferRepList.Bytes())
	if err != nil {
		return nil, err
	}

	err = ipZIPWriter.Close()
	if err != nil {
		return nil, err
	}

	return &ipZIPFile, nil
}

func cleansing(masukan string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				masukan, ",", "-"),
			"\"", ""),
		"'", "")
}
