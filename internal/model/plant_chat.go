package model

// PlantInfo 식물의 기본 정보를 담는 구조체
type PlantInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	// Environment      string `json:"environment"`
	// Characteristics  string `json:"characteristics"`
	// CareInstructions string `json:"careInstructions"`
	// Health           string `json:"health"`
	// Temperature      string `json:"temperature"`
	// Humidity         string `json:"humidity"`
}

// Message 대화 메시지를 담는 구조체
type Message struct {
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Time    string `json:"time,omitempty"` // 메시지 생성 시간
}

// PlantConditions 식물의 일일 상태 정보를 담는 구조체
type PlantConditions struct {
	Temperature          float64 `json:"temperature"`
	Humidity             float64 `json:"humidity"`
	PH                   float64 `json:"ph"`
	Rainfall             float64 `json:"rainfall"`
	SoilMoisture         float64 `json:"soil_moisture"`
	SunlightExposure     float64 `json:"sunlight_exposure"`
	WaterUsageEfficiency float64 `json:"water_usage_efficiency"`
	N                    float64 `json:"N"`
	P                    float64 `json:"P"`
	K                    float64 `json:"K"`
	SoilType             float64 `json:"soil_type"`
	WindSpeed            float64 `json:"wind_speed"`
	CO2Concentration     float64 `json:"co2_concentration"`
	CropDensity          float64 `json:"crop_density"`
	PestPressure         float64 `json:"pest_pressure"`
	UrbanAreaProximity   float64 `json:"urban_area_proximity"`
	FrostRisk            float64 `json:"frost_risk"`
}

// DailyPlantData 일일 식물 데이터를 담는 구조체
type DailyPlantData struct {
	Day        int             `json:"day"`
	Conditions PlantConditions `json:"conditions"`
}

// PlantChatRequest 식물 채팅 요청을 담는 구조체
type PlantChatRequest struct {
	UserID              string         `json:"user_id"`
	PlantInfo           PlantInfo      `json:"plant_info"`
	UserMessage         string         `json:"user_message"`
	ChildProfile        ChildProfile   `json:"child_profile"`
	DailyData           DailyPlantData `json:"daily_data"`
	ConversationHistory []Message      `json:"conversationHistory,omitempty"`
}

// PlantChatResponse 식물 채팅 응답을 담는 구조체
type PlantChatResponse struct {
	Message string `json:"message"`
}
