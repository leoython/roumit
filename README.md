在日常操作中，有大量如此操作：枚举一个id数组，并行对每个id做某种增删改除。为了避免因为某个id的不稳定操作(连接数据库，网络异常等)导致的整体效率下降，最常见的写法就是每个元素开一个独立的goroutine去处理。但是这种做法容易造成超高的并发导致拥堵，最终降低了整体效率。

这个新增的roumit库就是专门解决这个问题。使用Map或者更好一点的MapWithTimeout，可以控制最高并发数量workNums，和每个子任务允许的操作时间timeout。一个具体例子长这样：
```
Ids := []int{1, 2, 3}
items := make([]*Item, len(Ids))
err := routine.MapWithTimeout(len(items), 5, time.Second, func(ctx context.Context, i int) {
    item := getItemByID(Ids[i])
    if ctx.Err() == nil {
        items[i] = item
    }
})
```