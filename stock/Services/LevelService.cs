using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Grpc.Core;
using Microsoft.Extensions.Logging;

namespace stock
{
    public class LevelService : Level.LevelBase
    {
        private readonly ILogger<LevelService> _logger;
        public LevelService(ILogger<LevelService> logger)
        {
            _logger = logger;
        }

        public override Task<LevelReply> GetLevel(LevelRequest request, ServerCallContext context)
        {
            return Task.FromResult(new LevelReply
            {
                Message = "Level " + new Random().Next(0, 100)
            });
        }
    }
}
