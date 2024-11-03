namespace NestNetDesktop.ViewModels;

public partial class PostViewModel : ViewModelBase
{
    private string _title;
    private string _body;
    private string _imagemd5;
    private string _id;

    public string Id
    {
        get => _id;
        set => SetProperty(ref _id, value);
    }
    
    public string Title
    {
        get => _title;
        set => SetProperty(ref _title, value);
    }

    public string Body
    {
        get => _body;
        set => SetProperty(ref _body, value);
    }

    public string Imagemd5
    {
        get => _imagemd5;
        set => SetProperty(ref _imagemd5, value);
    }
}