package main

func createPost(post *Post ) error {
	_, err := db.Exec("INSERT INTO public.posts(title, content) VALUES ($1, $2);",
		post.Title,
		post.Content,
	)

	return err
}

func getPost(id string) (Post, error) {
	var p Post
	row := db.QueryRow("SELECT id, title, content, published, created_at from public.posts WHERE id=$1;",
		id,
	)

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Published, &p.CreatedAt)

	if err != nil {
		return Post{}, err 
	}

	return p, nil
}

func getPosts() ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, published from public.posts WHERE published = true;")

	if err != nil {
		return nil, err
	}

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Published)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func getDraft() ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, published from public.posts WHERE published = false;")

	if err != nil {
		return nil, err
	}

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Published)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func updatePost(id string, post *Post ) error {
	_, err := db.Exec("UPDATE public.posts SET title=$1, content=$2 WHERE id =$3;",
		post.Title,
		post.Content,
		id,
	)

	return err
}

func publishedPost(id string, post *Post ) error {
	_, err := db.Exec("UPDATE public.posts SET published=$1 WHERE id =$2;", 
		post.Published,
		id,
	)

	return err
}

func deletePost(id string) error {
	_, err := db.Exec("DELETE from public.posts WHERE id=$1;",
	id,
	)

	return err
}