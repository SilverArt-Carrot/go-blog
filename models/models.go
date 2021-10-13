package models

//title = models.CharField(verbose_name='标题', max_length=80)
//author = models.ForeignKey(User, verbose_name='作者', on_delete=models.CASCADE)
//body = models.TextField(verbose_name='正文')
//category = models.ForeignKey(Category, verbose_name='分类', on_delete=models.CASCADE)
//tag = models.ManyToManyField(Tag, verbose_name='标签', blank=True)
//img = models.ImageField(verbose_name='图片', default=None)
//excerpt = models.CharField(verbose_name='摘要', max_length=200, blank=True)
//created_time = models.DateTimeField(verbose_name='创建时间', default=timezone.now)
//modified_time = models.DateTimeField(verbose_name='修改时间')
//visits = models.PositiveIntegerField(verbose_name='访问量', default=0, editable=False)

type Time struct {
	Year int
	Month int
	Day int
}

//定义一个文章结构体
type Post struct {
	Id              	int     //primary key
	Title           	string
	Author          	string
	Body            	string
	Category        	int  	//foreign key,reference Id of Category
	Tag             	int	 	//foreign key,reference Id of Tag
	Img             	string  //图片所在路径
	Excerpt         	string
	CreatedTime 		Time
	ModifyTime 			Time
}


//name = models.CharField(max_length=50)

//定义一个文章分类结构体
type Category struct {
	Id   int    // pk
	Name string
}

//name = models.CharField(max_length=20)

//定义个一个文章标签结构体
type Tag struct {
	Id   int    // pk
	Name string
}

//定义一个归档标签结构体
type Archives struct {
	Year int
	Month int
}


