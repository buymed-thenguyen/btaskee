package db

import (
	dbModel "btaskee/model/db"
)

func SeedQuizzes() error {
	type RawQuiz struct {
		Title     string
		Questions []struct {
			Text    string
			Options []string // len = 4
			Answer  int      // 0-based index of correct answer
		}
	}

	quizzes := []RawQuiz{
		{
			Title: "Kiến thức tổng hợp",
			Questions: []struct {
				Text    string
				Options []string
				Answer  int
			}{
				{"Thủ đô của Việt Nam là gì?", []string{"Hà Nội", "TP.HCM", "Đà Nẵng", "Huế"}, 0},
				{"Tác giả Truyện Kiều là ai?", []string{"Nguyễn Du", "Hồ Xuân Hương", "Xuân Diệu", "Tố Hữu"}, 0},
				{"Đơn vị đo điện áp?", []string{"Volt", "Ampe", "Ohm", "Watt"}, 0},
				{"Sông dài nhất Việt Nam?", []string{"Sông Cửu Long", "Sông Hồng", "Sông Đồng Nai", "Sông Sài Gòn"}, 0},
				{"Ngày Quốc khánh Việt Nam?", []string{"2/9", "30/4", "1/5", "1/1"}, 0},
				{"Nguyên tố hóa học có ký hiệu H?", []string{"Hydro", "Heli", "Oxy", "Nito"}, 0},
				{"Động vật nào là vua rừng xanh?", []string{"Sư tử", "Hổ", "Báo", "Sói"}, 0},
				{"Màu cơ bản gồm?", []string{"Đỏ, xanh, vàng", "Tím, xanh lá, cam", "Xanh dương, đen, trắng", "Hồng, tím, nâu"}, 0},
				{"Ai là người đầu tiên lên mặt trăng?", []string{"Neil Armstrong", "Buzz Aldrin", "Yuri Gagarin", "Michael Collins"}, 0},
				{"Quốc gia có dân số cao nhất?", []string{"Trung Quốc", "Ấn Độ", "Mỹ", "Indonesia"}, 0},
			},
		},
		{
			Title:     "Toán học cơ bản",
			Questions: genMathQuestions(),
		},
		{
			Title:     "Địa lý Việt Nam",
			Questions: genGeoQuestions(),
		},
		{
			Title:     "Lịch sử thế giới",
			Questions: genHistoryQuestions(),
		},
		{
			Title:     "Khoa học tự nhiên",
			Questions: genScienceQuestions(),
		},
	}

	for _, raw := range quizzes {
		quiz := dbModel.Quiz{Title: raw.Title}
		if err := DB.Create(&quiz).Error; err != nil {
			return err
		}
		for _, q := range raw.Questions {
			question := dbModel.Question{QuizID: quiz.ID, QuestionText: q.Text, Score: 1}
			if err := DB.Create(&question).Error; err != nil {
				return err
			}
			for i, opt := range q.Options {
				ans := dbModel.AnswerOption{
					QuestionID: question.ID,
					Text:       opt,
					IsCorrect:  i == q.Answer,
				}
				if err := DB.Create(&ans).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func genMathQuestions() []struct {
	Text    string
	Options []string
	Answer  int
} {
	return []struct {
		Text    string
		Options []string
		Answer  int
	}{
		{"1 + 1 = ?", []string{"2", "3", "1", "0"}, 0},
		{"5 * 6 = ?", []string{"30", "11", "56", "26"}, 0},
		{"10 / 2 = ?", []string{"5", "2", "8", "10"}, 0},
		{"Căn bậc hai của 25 là?", []string{"5", "4", "6", "3"}, 0},
		{"20% của 50 là?", []string{"10", "5", "15", "20"}, 0},
		{"Số nguyên tố đầu tiên?", []string{"2", "3", "1", "5"}, 0},
		{"10 + 15 = ?", []string{"25", "20", "15", "30"}, 0},
		{"50 - 18 = ?", []string{"32", "42", "28", "38"}, 0},
		{"9 x 9 = ?", []string{"81", "72", "91", "99"}, 0},
		{"7^2 = ?", []string{"49", "36", "64", "25"}, 0},
	}
}

func genGeoQuestions() []struct {
	Text    string
	Options []string
	Answer  int
} {
	return []struct {
		Text    string
		Options []string
		Answer  int
	}{
		{"Tỉnh cực Bắc Việt Nam?", []string{"Hà Giang", "Lào Cai", "Cao Bằng", "Lạng Sơn"}, 0},
		{"Đảo lớn nhất Việt Nam?", []string{"Phú Quốc", "Côn Đảo", "Cát Bà", "Lý Sơn"}, 0},
		{"Thành phố trực thuộc trung ương?", []string{"Đà Nẵng", "Cần Thơ", "Huế", "Nha Trang"}, 0},
		{"Tỉnh nào thuộc miền Trung?", []string{"Nghệ An", "Bắc Ninh", "Hưng Yên", "Thái Bình"}, 0},
		{"Sông lớn nhất miền Nam?", []string{"Cửu Long", "Sài Gòn", "Đồng Nai", "Tiền Giang"}, 0},
		{"Tỉnh ven biển miền Trung?", []string{"Quảng Nam", "Kon Tum", "Gia Lai", "Đắk Lắk"}, 0},
		{"Biển nào giáp Việt Nam?", []string{"Biển Đông", "Biển Đỏ", "Biển Đen", "Biển Bắc"}, 0},
		{"Tỉnh có núi Fansipan?", []string{"Lào Cai", "Yên Bái", "Hà Giang", "Sơn La"}, 0},
		{"Biển lớn nhất Việt Nam?", []string{"Vịnh Bắc Bộ", "Vịnh Thái Lan", "Vịnh Hạ Long", "Vịnh Cam Ranh"}, 0},
		{"Tỉnh nào không giáp biển?", []string{"Hà Nội", "Quảng Ngãi", "Khánh Hòa", "Bình Thuận"}, 0},
	}
}

func genHistoryQuestions() []struct {
	Text    string
	Options []string
	Answer  int
} {
	return []struct {
		Text    string
		Options []string
		Answer  int
	}{
		{"Chiến tranh thế giới thứ hai kết thúc năm nào?", []string{"1945", "1940", "1954", "1939"}, 0},
		{"Quốc gia nào bị thả bom nguyên tử?", []string{"Nhật Bản", "Đức", "Liên Xô", "Anh"}, 0},
		{"Người sáng lập chủ nghĩa cộng sản?", []string{"Karl Marx", "Lenin", "Engels", "Stalin"}, 0},
		{"Tổng thống Mỹ trong WWII?", []string{"Franklin D. Roosevelt", "John F. Kennedy", "Bush", "Truman"}, 0},
		{"Nhà nước đầu tiên ở Việt Nam?", []string{"Văn Lang", "Âu Lạc", "Lý", "Trần"}, 0},
		{"Cách mạng tháng 8 diễn ra năm nào?", []string{"1945", "1954", "1930", "1975"}, 0},
		{"Chiến dịch Điện Biên Phủ?", []string{"1954", "1946", "1975", "1968"}, 0},
		{"Mỹ rút quân khỏi VN?", []string{"1973", "1975", "1968", "1970"}, 0},
		{"Vua Quang Trung tên thật?", []string{"Nguyễn Huệ", "Nguyễn Ánh", "Lê Lợi", "Lê Lai"}, 0},
		{"Thời kỳ đồ đá bắt đầu khi nào?", []string{"2 triệu năm TCN", "1 triệu năm TCN", "10.000 năm TCN", "5.000 năm TCN"}, 0},
	}
}

func genScienceQuestions() []struct {
	Text    string
	Options []string
	Answer  int
} {
	return []struct {
		Text    string
		Options []string
		Answer  int
	}{
		{"Nước sôi ở bao nhiêu độ C?", []string{"100", "90", "80", "120"}, 0},
		{"Cơ quan hô hấp của người?", []string{"Phổi", "Gan", "Tim", "Dạ dày"}, 0},
		{"Tốc độ ánh sáng gần đúng?", []string{"300,000 km/s", "150,000 km/s", "100,000 km/s", "1 triệu km/s"}, 0},
		{"Mắt người nhìn thấy ánh sáng nào?", []string{"Ánh sáng khả kiến", "Hồng ngoại", "Tia X", "Tia gamma"}, 0},
		{"Lực hút của Trái Đất gọi là gì?", []string{"Trọng lực", "Ma sát", "Từ lực", "Động lực"}, 0},
		{"Máu đỏ do gì tạo ra?", []string{"Hồng cầu", "Bạch cầu", "Tiểu cầu", "Tủy xương"}, 0},
		{"Thực vật cần gì để quang hợp?", []string{"Ánh sáng", "Oxy", "Nhiệt độ", "Gió"}, 0},
		{"Âm thanh lan truyền qua?", []string{"Không khí", "Chân không", "Ánh sáng", "Điện"}, 0},
		{"Hệ mặt trời có mấy hành tinh?", []string{"8", "9", "7", "10"}, 0},
		{"Trái đất quay quanh gì?", []string{"Mặt trời", "Mặt trăng", "Sao Hỏa", "Trục"}, 0},
	}
}
