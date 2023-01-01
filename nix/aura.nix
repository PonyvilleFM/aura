self:
{ config, lib, pkgs, ... }:
let cfg = config.within.services.aura;
in
with lib; {
  options.within.services.aura = {
    enable = mkEnableOption "Activates Aura (the PonyvilleFM chatbot)";

    port = mkOption {
      type = types.port;
      default = 58238;
      example = 9001;
      description = "The port number mi should listen on for HTTP traffic";
    };

    domain = mkOption {
      type = types.str;
      default = "pvfmsets.cf";
      example = "pvfm.akua";
      description =
        "The domain name that nginx should check against for HTTP hostnames";
    };

    wwwRedir = mkEnableOption "Redirect www. to the domain root";
  };

  config = mkIf cfg.enable {
    users.users.aura = {
      createHome = true;
      description = "github.com/PonyvilleFM/aura";
      isSystemUser = true;
      group = "within";
      home = "/srv/within/aura";
      extraGroups = [ "keys" ];
    };

    systemd.services.aura = {
      wantedBy = [ "multi-user.target" ];
      after = [ "aura-key.service" ];
      wants = [ "aura-key.service" ];

      serviceConfig = {
        User = "aura";
        Group = "within";
        Restart = "on-failure";
        WorkingDirectory = "/srv/within/aura";
        RestartSec = "30s";
      };

      script = let aura = self.packages.${pkgs.system}.default;
      in ''
        export PATH=$PATH:${pkgs.streamripper}/bin
        export PORT=${toString cfg.port}
        export RECORDING_DOMAIN=${cfg.domain}
        env
        exec ${aura}/bin/aura
      '';
    };

    services.nginx.virtualHosts."aura" = {
      serverName = "${cfg.domain}";
      forceSSL = true;
      useACMEHost = "pvfmsets.cf";

      locations."/".proxyPass = "http://127.0.0.1:${toString cfg.port}";
    };
  };
}
