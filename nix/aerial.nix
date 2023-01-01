self:
{ config, lib, pkgs, ... }:
with lib; {
  options.within.services.aerial.enable =
    mkEnableOption "Activates Aerial (the PonyvilleFM chatbot)";

  config = mkIf config.within.services.aerial.enable {
    users.users.aerial = {
      createHome = true;
      description = "github.com/PonyvilleFM/aura";
      isSystemUser = true;
      group = "within";
      home = "/srv/within/aerial";
      extraGroups = [ "keys" ];
    };

    systemd.services.aerial = {
      wantedBy = [ "multi-user.target" ];
      after = [ "aerial-key.service" ];
      wants = [ "aerial-key.service" ];

      serviceConfig = {
        User = "aerial";
        Group = "within";
        Restart = "on-failure";
        WorkingDirectory = "/srv/within/aerial";
        RestartSec = "30s";
      };

      script = let aura = self.packages.${pkgs.system}.default;
      in ''
        exec ${aura}/bin/aerial
      '';
    };
  };
}
