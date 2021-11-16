package v1

//TODO已完成: 获取现场确认时间段数据列表
//TODO已完成: 根据ID删除现场确认时间段数据(管理员)
//TODO已完成: 根据ID修改现场确认时间段数据(管理员)
//TODO已完成: 新增现场确认时间段数据(管理员)

type NewTimeBody struct {
	Name        string `json:"timeName" form:"timeName" validate:"required"`       //时间段
	OptionalNum uint   `json:"optionalNum" form:"optionalNum" validate:"required"` //可选数
	StagingId   uint   `json:"staging_id" form:"staging_id" validate:"required"`   //分期id
	//SelectedNum uint   `json:"selectedNum" form:"selectedNum" validate:"ltefield=OptionalNum"` //已选数
}

type GetTimeByStagingId struct {
	StagingId uint `json:"staging_id" form:"staging_id" validate:"required"`
}

type UpdateTimeBody struct {
	Id          uint   `json:"id" form:"id" validate:"required"`
	Name        string `json:"name" form:"name" `                //时间段
	OptionalNum uint   `json:"optional_num" form:"optional_num"` //可选数
	//	SelectedNum uint   `json:"selected_num" form:"selected_num" ` //已选数
}

type DeleteTimeBody struct {
	Id uint `json:"id" form:"id" validate:"required"`
}
