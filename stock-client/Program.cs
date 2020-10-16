using System;
using System.Threading.Tasks;
using Grpc.Net.Client;

namespace stock_client
{
    class Program
    {
        static async Task Main(string[] args)
        {
            // The port number(5001) must match the port of the gRPC server.
            using var channel = GrpcChannel.ForAddress("https://localhost:5001");
            var client = new Level.LevelClient(channel);
            var reply = client.GetLevel(
                              new LevelRequest { Name = "LevelClient" });
            Console.WriteLine("Level: " + reply.Message);
            Console.WriteLine("Press any key to exit...");
            Console.ReadKey();
        }
    }
}
