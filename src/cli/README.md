```bash
todo register "username" "email" "password"
todo login "email" "password"
todo logout

todo create board "title"
todo create column *board_idx* "title"
todo create card *board_idx* *column_idx* "title" ["description"]

todo show boards
todo show board *board_idx*
todo show cards *board_idx* *column_idx*
todo show card *board_idx* *column_idx* *card_idx*

todo update board *board_idx* title "title"
todo update column *board_idx* *column_idx* title "title"
todo update card *board_idx* *column_idx* *card_idx* title "title"
todo update card *board_idx* *column_idx* *card_idx* description "description"

todo delete board *board_idx*
todo delete column *board_idx* *column_idx*
todo delete card *board_idx* *column_idx* *card_idx*

todo stats from "DD-MM-YYYY" to "DD-MM-YYYY"
```
