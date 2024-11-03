using System.Collections.ObjectModel;
using System.Net.Http;
using System.Net.Http.Json;
using System.Security.Authentication.ExtendedProtection;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using CommunityToolkit.Mvvm.Input;
using NestNetDesktop.Models;

namespace NestNetDesktop.ViewModels;

public partial class MainWindowViewModel : ViewModelBase
{
    private string _instanceUrl;
    private string _userName;
    private ObservableCollection<PostViewModel> _posts;
    public string Greeting => "Welcome to Avalonia!";


    public MainWindowViewModel() : base()
    {
        InstanceUrl = "http://localhost:4929";
        NewPost = new NewPostViewModel(InstanceUrl);
        Posts = new ObservableCollection<PostViewModel>(); 
    }

    public string InstanceUrl
    {
        get => _instanceUrl;
        set
        {
            SetProperty(ref _instanceUrl, value);
            if (NewPost is not null)
            {
                NewPost.InstanceUrl = value;
            }
            
        }
    }

    public string UserName
    {
        get => _userName;
        set => SetProperty(ref _userName, value);
    }

    public ObservableCollection<PostViewModel> Posts
    {
        get => _posts;
        set => SetProperty(ref _posts, value);
    }

    public NewPostViewModel NewPost { get; }

    [RelayCommand]
    private async Task PostSetName()
    {
        var client = new HttpClient();
        var name = new NameModel(UserName);
        var body = new StringContent(JsonSerializer.Serialize(name), Encoding.UTF8, "application/json");
        var response = await client.PostAsync($"{InstanceUrl}/set_name", body);
    }

    [RelayCommand]
    private async Task Refresh()
    {
        var client = new HttpClient();
        // Get currently set name /get_name and set UserName to it
        var name = await (await client.GetAsync($"{InstanceUrl}/get_name")).Content.ReadFromJsonAsync<NameModel>();
        // Get all posts with /retrieve
        var posts = await client.GetFromJsonAsync<PostViewModel[]>($"{InstanceUrl}/retrieve");
        Posts = new ObservableCollection<PostViewModel>(posts);
    }
}