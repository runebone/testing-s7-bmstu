```bash
todo register "username" "email" "password"
todo login "email" "password"
todo logout

todo create board "title"
todo create column *board_id* "title"
todo create card *board_id* *column_id* "title" ["description"]

todo show boards
todo show board *board_id*
todo show cards *board_id* *column_id*
todo show card *board_id* *column_id* *card_id*

todo update board *board_id* title "title"
todo update column *board_id* *column_id* title "title"
todo update card *board_id* *column_id* *card_id* title "title"
todo update card *board_id* *column_id* *card_id* description "description"

todo delete board *board_id*
todo delete column *board_id* *column_id*
todo delete card *board_id* *column_id* *card_id*

todo stats from "DD-MM-YYYY" to "DD-MM-YYYY"
```
