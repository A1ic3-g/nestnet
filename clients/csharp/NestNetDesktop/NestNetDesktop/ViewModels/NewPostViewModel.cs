using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using CommunityToolkit.Mvvm.Input;

namespace NestNetDesktop.ViewModels;

public partial class NewPostViewModel : PostViewModel
{
    [JsonIgnore] public string InstanceUrl;

    public NewPostViewModel(string instanceUrl) : base()
    {
        InstanceUrl = instanceUrl;
    }

    [RelayCommand]
    private async Task Publish()
    {
        var client = new HttpClient();
        var body = JsonSerializer.Serialize(this);
        await client.PostAsync($"{InstanceUrl}/add_post",
            new StringContent(body, Encoding.UTF8, "application/json"));
    }
}