package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sample Go App", func() {
	var (
		fromBrokerOneTag           message
		fromBrokerOneTagBody       []messageBody
		fromBrokerOneTagBody0      messageBody
		fromBrokerOneTagDatapoints [][]interface{}

		fromBrokerTwoTags                message
		fromBrokerTwoTagsBody            []messageBody
		fromBrokerTwoTagsBody0           messageBody
		fromBrokerTwoTagsBody1           messageBody
		fromBrokerTwoTagsBody0Datapoints [][]interface{}
		fromBrokerTwoTagsBody1Datapoints [][]interface{}

		fromBrokerTwoTagsSame                message
		fromBrokerTwoTagsSameBody            []messageBody
		fromBrokerTwoTagsSameBody0           messageBody
		fromBrokerTwoTagsSameBody1           messageBody
		fromBrokerTwoTagsSameBody0Datapoints [][]interface{}
		fromBrokerTwoTagsSameBody1Datapoints [][]interface{}

		oneTagNoChange           message
		oneTagNoChangeBody       []messageBody
		oneTagNoChangeBody0      messageBody
		oneTagNoChangeDatapoints [][]interface{}

		oneTagYesChange           message
		oneTagYesChangeBody       []messageBody
		oneTagYesChangeBody0      messageBody
		oneTagYesChangeDatapoints [][]interface{}

		twoTagsNoChange                message
		twoTagsNoChangeBody            []messageBody
		twoTagsNoChangeBody0           messageBody
		twoTagsNoChangeBody1           messageBody
		twoTagsNoChangeBody0Datapoints [][]interface{}
		twoTagsNoChangeBody1Datapoints [][]interface{}

		twoTagsOneChange                message
		twoTagsOneChangeBody            []messageBody
		twoTagsOneChangeBody0           messageBody
		twoTagsOneChangeBody1           messageBody
		twoTagsOneChangeBody0Datapoints [][]interface{}
		twoTagsOneChangeBody1Datapoints [][]interface{}

		twoTagsBothChange                message
		twoTagsBothChangeBody            []messageBody
		twoTagsBothChangeBody0           messageBody
		twoTagsBothChangeBody1           messageBody
		twoTagsBothChangeBody0Datapoints [][]interface{}
		twoTagsBothChangeBody1Datapoints [][]interface{}
	)

	Describe("Behavior of the scaleData function", func() {
		var attributes = make(map[string]string)

		Context("When only one tag is in the input", func() {
			//Construct JSON in the form
			// from_broker_one_tag = {
			// 	body: [
			// 	{
			// 		attributes:{
			// 			machine_type:"opcua"
			// 		},
			// 		datapoints: [[1537377630622,80.0,3]],
			// 		name: "My.App.DOUBLE1"
			// 	}
			// 	],
			// 	messageId: "flex-pipe"
			// };
			fromBrokerOneTagDatapoint := []interface{}{1537377630622, 80.0, 3}
			fromBrokerOneTagDatapoints := append(fromBrokerOneTagDatapoints, fromBrokerOneTagDatapoint)
			attributes["machine_type"] = "opcua"
			fromBrokerOneTagBody0 = messageBody{
				Attributes: attributes,
				Datapoints: fromBrokerOneTagDatapoints,
				Name:       "My.App.DOUBLE1",
			}
			fromBrokerOneTagBody = append(fromBrokerOneTagBody, fromBrokerOneTagBody0)
			fromBrokerOneTag = message{
				MessageID: "flex-pipe",
				Body:      fromBrokerOneTagBody,
			}
			attributes["machine_type"] = "opcua"
			It("if the tagToMatch does not match the tag on the data to scale, it should not be scaled.", func() {
				//Construct JSON in the form
				// one_tag_no_change = {
				// 	body: [
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints: [[1537377630622,80.0,3]],
				// 		name: "My.App.DOUBLE1"
				// 	}
				// 	],
				// 	messageId: "flex-pipe"
				// };
				oneTagNoChangeDatapoint := []interface{}{1537377630622, 80.0, 3}
				oneTagNoChangeDatapoints := append(oneTagNoChangeDatapoints, oneTagNoChangeDatapoint)

				oneTagNoChangeBody0 = messageBody{
					Attributes: attributes,
					Datapoints: oneTagNoChangeDatapoints,
					Name:       "My.App.DOUBLE1",
				}
				oneTagNoChangeBody = append(oneTagNoChangeBody, oneTagNoChangeBody0)
				oneTagNoChange = message{
					MessageID: "flex-pipe",
					Body:      oneTagNoChangeBody,
				}
				Expect(scaleData(fromBrokerOneTag, "My.App.DOUBLE2")).To(Equal(oneTagNoChange))
			})
			It("if the tagToMatch matches the tag on the data to scale, the data should be scaled by 1000, and the tag name changed.", func() {
				//Construct JSON in the form
				// one_tag_yes_change = {
				// 	body: [
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints: [[1537377630622,80000,3]],
				// 		name: "My.App.DOUBLE1.scaled_by_1000"
				// 	}
				// 	],
				// 	messageId: "flex-pipe"
				// };
				oneTagYesChangeDatapoint := []interface{}{1537377630622, 80000.0, 3}
				oneTagYesChangeDatapoints := append(oneTagYesChangeDatapoints, oneTagYesChangeDatapoint)

				oneTagYesChangeBody0 = messageBody{
					Attributes: attributes,
					Datapoints: oneTagYesChangeDatapoints,
					Name:       "My.App.DOUBLE1.scaled_by_1000",
				}
				oneTagYesChangeBody = append(oneTagYesChangeBody, oneTagYesChangeBody0)
				oneTagYesChange = message{
					MessageID: "flex-pipe",
					Body:      oneTagYesChangeBody,
				}
				Expect(scaleData(fromBrokerOneTag, "My.App.DOUBLE1")).To(Equal(oneTagYesChange))

			})

		})
		Context("When there are two tags in the input", func() {
			//Construct JSON in the form
			// from_broker_two_tags = {
			// 	body: [
			// 	{
			// 		attributes:{
			// 			machine_type:"opcua"
			// 		},
			// 		datapoints: [[1537377630622,80.0,3]],
			// 		name: "My.App.DOUBLE1"
			// 	},
			// 	{
			// 		attributes:{
			// 			machine_type:"opcua"
			// 		},
			// 		datapoints:[[1537377630622,112.64,3]],
			// 		name: "My.App.DOUBLE2"
			// 	}
			// 	],
			// 	messageId: "flex-pipe"
			// };
			fromBrokerTwoTagsBody0Datapoint := []interface{}{1537377630622, 80.0, 3}

			fromBrokerTwoTagsBody1Datapoint := []interface{}{1537377630622, 112.64, 3}
			fromBrokerTwoTagsBody0Datapoints := append(fromBrokerTwoTagsBody0Datapoints, fromBrokerTwoTagsBody0Datapoint)
			fromBrokerTwoTagsBody1Datapoints := append(fromBrokerTwoTagsBody1Datapoints, fromBrokerTwoTagsBody1Datapoint)
			fromBrokerTwoTagsBody0 = messageBody{
				Attributes: attributes,
				Datapoints: fromBrokerTwoTagsBody0Datapoints,
				Name:       "My.App.DOUBLE1",
			}
			fromBrokerTwoTagsBody = append(fromBrokerTwoTagsBody, fromBrokerTwoTagsBody0)
			fromBrokerTwoTagsBody1 = messageBody{
				Attributes: attributes,
				Datapoints: fromBrokerTwoTagsBody1Datapoints,
				Name:       "My.App.DOUBLE2",
			}
			fromBrokerTwoTagsBody = append(fromBrokerTwoTagsBody, fromBrokerTwoTagsBody1)

			fromBrokerTwoTags = message{
				MessageID: "flex-pipe",
				Body:      fromBrokerTwoTagsBody,
			}
			It("if the tagToMatch does not match the tags on the data to scale, none should not be scaled.", func() {
				//Construct JSON in the form
				// two_tags_no_change = {
				// 	body: [
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints: [[1537377630622,80.0,3]],
				// 		name: "My.App.DOUBLE1"
				// 	},
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints:[[1537377630622,112.64,3]],
				// 		name: "My.App.DOUBLE2"
				// 	}
				// 	],
				// 	messageId: "flex-pipe"
				// };
				twoTagsNoChangeBody0Datapoint := []interface{}{1537377630622, 80.0, 3}

				twoTagsNoChangeBody1Datapoint := []interface{}{1537377630622, 112.64, 3}
				twoTagsNoChangeBody0Datapoints := append(twoTagsNoChangeBody0Datapoints, twoTagsNoChangeBody0Datapoint)
				twoTagsNoChangeBody1Datapoints := append(twoTagsNoChangeBody1Datapoints, twoTagsNoChangeBody1Datapoint)
				twoTagsNoChangeBody0 = messageBody{
					Attributes: attributes,
					Datapoints: twoTagsNoChangeBody0Datapoints,
					Name:       "My.App.DOUBLE1",
				}
				twoTagsNoChangeBody = append(twoTagsNoChangeBody, twoTagsNoChangeBody0)
				twoTagsNoChangeBody1 = messageBody{
					Attributes: attributes,
					Datapoints: twoTagsNoChangeBody1Datapoints,
					Name:       "My.App.DOUBLE2",
				}
				twoTagsNoChangeBody = append(twoTagsNoChangeBody, twoTagsNoChangeBody1)

				twoTagsNoChange = message{
					MessageID: "flex-pipe",
					Body:      twoTagsNoChangeBody,
				}
				Expect(scaleData(fromBrokerTwoTags, "My.App.DOUBLE3")).To(Equal(twoTagsNoChange))

			})
			It("if the tagToMatch does matches one tag on the data to scale, the appropriate data should be scaled.", func() {
				//Construct JSON in the form
				// two_tags_one_change = {
				// 	body: [
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints: [[1537377630622,80000,3]],
				// 		name: "My.App.DOUBLE1.scaled_by_1000"
				// 	},
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints:[[1537377630622,112.64,3]],
				// 		name: "My.App.DOUBLE2"
				// 	}
				// 	],
				// 	messageId: "flex-pipe"
				// };
				twoTagsOneChangeBody0Datapoint := []interface{}{1537377630622, 80000.0, 3}

				twoTagsOneChangeBody1Datapoint := []interface{}{1537377630622, 112.64, 3}
				twoTagsOneChangeBody0Datapoints := append(twoTagsOneChangeBody0Datapoints, twoTagsOneChangeBody0Datapoint)
				twoTagsOneChangeBody1Datapoints := append(twoTagsOneChangeBody1Datapoints, twoTagsOneChangeBody1Datapoint)
				twoTagsOneChangeBody0 = messageBody{
					Attributes: attributes,
					Datapoints: twoTagsOneChangeBody0Datapoints,
					Name:       "My.App.DOUBLE1.scaled_by_1000",
				}
				twoTagsOneChangeBody = append(twoTagsOneChangeBody, twoTagsOneChangeBody0)
				twoTagsOneChangeBody1 = messageBody{
					Attributes: attributes,
					Datapoints: twoTagsOneChangeBody1Datapoints,
					Name:       "My.App.DOUBLE2",
				}
				twoTagsOneChangeBody = append(twoTagsOneChangeBody, twoTagsOneChangeBody1)

				twoTagsOneChange = message{
					MessageID: "flex-pipe",
					Body:      twoTagsOneChangeBody,
				}
				Expect(scaleData(fromBrokerTwoTags, "My.App.DOUBLE1")).To(Equal(twoTagsOneChange))

			})
			It("if the tagToMatch does matches both tags on the data to scale, the appropriate data should be scaled.", func() {
				//Construct JSON in the form
				// from_broker_two_tags_same = {
				// 	body: [
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints: [[1537377630622,80.0,3]],
				// 		name: "My.App.DOUBLE1"
				// 	},
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints:[[1537377630622,112.64,3]],
				// 		name: "My.App.DOUBLE1"
				// 	}
				// 	],
				// 	messageId: "flex-pipe"
				// };
				fromBrokerTwoTagsSameBody0Datapoint := []interface{}{1537377630622, 80.0, 3}

				fromBrokerTwoTagsSameBody1Datapoint := []interface{}{1537377630622, 112.64, 3}
				fromBrokerTwoTagsSameBody0Datapoints := append(fromBrokerTwoTagsSameBody0Datapoints, fromBrokerTwoTagsSameBody0Datapoint)
				fromBrokerTwoTagsSameBody1Datapoints := append(fromBrokerTwoTagsSameBody1Datapoints, fromBrokerTwoTagsSameBody1Datapoint)
				fromBrokerTwoTagsSameBody0 = messageBody{
					Attributes: attributes,
					Datapoints: fromBrokerTwoTagsSameBody0Datapoints,
					Name:       "My.App.DOUBLE1",
				}
				fromBrokerTwoTagsSameBody = append(fromBrokerTwoTagsSameBody, fromBrokerTwoTagsSameBody0)
				fromBrokerTwoTagsSameBody1 = messageBody{
					Attributes: attributes,
					Datapoints: fromBrokerTwoTagsSameBody1Datapoints,
					Name:       "My.App.DOUBLE1",
				}
				fromBrokerTwoTagsSameBody = append(fromBrokerTwoTagsSameBody, fromBrokerTwoTagsSameBody1)

				fromBrokerTwoTagsSame = message{
					MessageID: "flex-pipe",
					Body:      fromBrokerTwoTagsSameBody,
				}
				//Construct JSON in the form
				// two_tags_both_change = {
				// 	body: [
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints: [[1537377630622,80000,3]],
				// 		name: "My.App.DOUBLE1.scaled_by_1000"
				// 	},
				// 	{
				// 		attributes:{
				// 			machine_type:"opcua"
				// 		},
				// 		datapoints:[[1537377630622,112640.0,3]],
				// 		name: "My.App.DOUBLE1.scaled_by_1000"
				// 	}
				// 	],
				// 	messageId: "flex-pipe"
				// };
				twoTagsBothChangeBody0Datapoint := []interface{}{1537377630622, 80000.0, 3}

				twoTagsBothChangeBody1Datapoint := []interface{}{1537377630622, 112640.0, 3}
				twoTagsBothChangeBody0Datapoints := append(twoTagsBothChangeBody0Datapoints, twoTagsBothChangeBody0Datapoint)
				twoTagsBothChangeBody1Datapoints := append(twoTagsBothChangeBody1Datapoints, twoTagsBothChangeBody1Datapoint)
				twoTagsBothChangeBody0 = messageBody{
					Attributes: attributes,
					Datapoints: twoTagsBothChangeBody0Datapoints,
					Name:       "My.App.DOUBLE1.scaled_by_1000",
				}
				twoTagsBothChangeBody = append(twoTagsBothChangeBody, twoTagsBothChangeBody0)
				twoTagsBothChangeBody1 = messageBody{
					Attributes: attributes,
					Datapoints: twoTagsBothChangeBody1Datapoints,
					Name:       "My.App.DOUBLE1.scaled_by_1000",
				}
				twoTagsBothChangeBody = append(twoTagsBothChangeBody, twoTagsBothChangeBody1)

				twoTagsBothChange = message{
					MessageID: "flex-pipe",
					Body:      twoTagsBothChangeBody,
				}
				Expect(scaleData(fromBrokerTwoTagsSame, "My.App.DOUBLE1")).To(Equal(twoTagsBothChange))

			})
		})
	})
})
