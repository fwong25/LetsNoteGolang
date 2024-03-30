# LetsNote App

### App Features
- Create notes with title and content
- Notes are editable
- Also shows last modified date

<img src="assets/interface1.png" alt="drawing" width="500"/>

### Prerequisites
- Backend: Golang, Postgresql
- Frontend: HTML, Bootstrap

This project is migrated from python+django version ([Link](https://github.com/fwong25/LetsNoteDjango))

Frontend UI design reference: [Todo App with Python+Django](https://www.youtube.com/watch?v=Nnoxz9JGdLU&ab_channel=CodAffection)

Before running project, make sure you have created the DB table in postgresql with following cmd
```
CREATE TABLE letsnote_note (
	id SERIAL PRIMARY KEY NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_date TEXT NOT NULL,
	last_modified_date TEXT NOT NULL
);
```

### Running the project
```
go run main.go
```

In your browser, e.g. Chrome, open the following address:
http://127.0.0.1:8000/list_note
or
http://localhost:8000/list_note

<img src="assets/interface2.png" alt="drawing" width="500"/>

### App Flow
<img src="assets/app_flow.png" alt="drawing" width="500"/>

### Database Design
<img src="assets/db_table.png" alt="drawing" width="500"/>