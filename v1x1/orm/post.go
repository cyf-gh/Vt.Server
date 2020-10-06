// TODO: 文章是否可以删除？
// TODO: 文章的访问权限问题？
package orm

import (
	err "../err"
)

// 上传的文章
type Post struct {
	Id int64
	Title string
	Text string
	TagIds[] int64
	OwnerId int64
	IsPrivate bool
}

// 上传的tag标签结构
// tag同时将作为分区的主要依据
// TODO:tag达到一定数量则升级为分区
type Tag struct {
	Id int64
	Text string
	IsCatalog bool
}

func Sync2Post() {
	e := engine_post.Sync2(new(Post))
	err.Check( e )
	e = engine_post.Sync2(new(Tag))
	err.Check( e )
}

// 通过某一个人获取所有他的文章
// 仅通过atk
func GetPostsByOwnerAll( OwnerId int64 ) ( []Post, error ) {
	var posts []Post
	e := engine_post.Table("Post").Where( "owner_id = ?", OwnerId).Find(&posts)
	return posts, e
}

// 通过某一个人获取所有公开文章
func GetPostsByOwnerPublic( OwnerId int64 ) ( []Post, error ) {
	var posts []Post
	e := engine_post.Table("Post").Where( "owner_id = ? and is_private = 1", OwnerId).Find(&posts)
	return posts, e
}

func GetPostsAll() ( []Post, error ) {
	var posts []Post
	e := engine_post.Table("Post").Where( "is_private = 1" ).Find(&posts)
	return posts, e
}

// 向数据库添加一笔新文章
func NewPost( title, text string, owner int64, tags []string) error {
	tagIds, e  := GetTagIds( tags )

	_, e = engine_post.Table("Post").Insert( &Post{
		Title:     title,
		Text:      text,
		TagIds:    tagIds,
		OwnerId: owner,
	})
	// err.Check( e )
	return e
}

// 修改文章
func ModifyPost( id int64, title, text string, owner int64, tags []string) error {
	tagIds, e := GetTagIds( tags )

	_, e = engine_post.Table("Post").ID(id).Update(&Post{
		Title:     title,
		Text:      text,
		TagIds:    tagIds,
		OwnerId:   owner,
	})
	return e
}

// 修改文章，不修改内容
// 减轻流量负担
func ModifyPostNoText( id int64, title string, owner int64, tags []string) error {
	tagIds, e := GetTagIds( tags )

	_, e = engine_post.Table("Post").ID(id).Update(&Post{
		Title:     title,
		TagIds:    tagIds,
		OwnerId:   owner,
	})
	return e
}

// 根据一系列tag的名字获取所有的tag的id
// 如果tag不存在则会被创建
func GetTagIds( tags []string ) ( []int64, error ) {
	var tagIds []int64

	for _, tagText := range tags  {
		t := new(Tag)
		GetTag:
		exists, _ := engine_post.Table("Tag").Where( "Text = ?", tagText).Get(t)

		// 如果当前没有该tag，则创建一个新tag
		if !exists {
			_, e := engine_post.Table("Tag").Insert( &Tag {
				Text:tagText,
				IsCatalog:false,
			})
			if e != nil {
				return nil, e
			}
			goto GetTag // 再次获取tag
		}
		// 这里t应该已经被填充
		tagIds = append(tagIds, t.Id)
	}
	return tagIds, nil
}