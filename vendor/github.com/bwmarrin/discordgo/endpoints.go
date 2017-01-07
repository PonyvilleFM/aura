// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains variables for all known Discord end points.  All functions
// throughout the Discordgo package use these variables for all connections
// to Discord.  These are all exported and you may modify them if needed.

package discordgo

// Known Discord API Endpoints.
var (
	EndpointStatus     = "https://status.discordapp.com/api/v2/"
	EndpointSm         = EndpointStatus + "scheduled-maintenances/"
	EndpointSmActive   = EndpointSm + "active.json"
	EndpointSmUpcoming = EndpointSm + "upcoming.json"

	EndpointDiscord  = "https://discordapp.com/"
	EndpointAPI      = EndpointDiscord + "api/"
	EndpointGuilds   = EndpointAPI + "guilds/"
	EndpointChannels = EndpointAPI + "channels/"
	EndpointUsers    = EndpointAPI + "users/"
	EndpointGateway  = EndpointAPI + "gateway"
	EndpointWebhooks = EndpointAPI + "webhooks/"

	EndpointAuth           = EndpointAPI + "auth/"
	EndpointLogin          = EndpointAuth + "login"
	EndpointLogout         = EndpointAuth + "logout"
	EndpointVerify         = EndpointAuth + "verify"
	EndpointVerifyResend   = EndpointAuth + "verify/resend"
	EndpointForgotPassword = EndpointAuth + "forgot"
	EndpointResetPassword  = EndpointAuth + "reset"
	EndpointRegister       = EndpointAuth + "register"

	EndpointVoice        = EndpointAPI + "/voice/"
	EndpointVoiceRegions = EndpointVoice + "regions"
	EndpointVoiceIce     = EndpointVoice + "ice"

	EndpointTutorial           = EndpointAPI + "tutorial/"
	EndpointTutorialIndicators = EndpointTutorial + "indicators"

	EndpointTrack        = EndpointAPI + "track"
	EndpointSso          = EndpointAPI + "sso"
	EndpointReport       = EndpointAPI + "report"
	EndpointIntegrations = EndpointAPI + "integrations"

	EndpointUser              = func(uID string) string { return EndpointUsers + uID }
	EndpointUserAvatar        = func(uID, aID string) string { return EndpointUsers + uID + "/avatars/" + aID + ".jpg" }
	EndpointUserSettings      = func(uID string) string { return EndpointUsers + uID + "/settings" }
	EndpointUserGuilds        = func(uID string) string { return EndpointUsers + uID + "/guilds" }
	EndpointUserGuild         = func(uID, gID string) string { return EndpointUsers + uID + "/guilds/" + gID }
	EndpointUserGuildSettings = func(uID, gID string) string { return EndpointUsers + uID + "/guilds/" + gID + "/settings" }
	EndpointUserChannels      = func(uID string) string { return EndpointUsers + uID + "/channels" }
	EndpointUserDevices       = func(uID string) string { return EndpointUsers + uID + "/devices" }
	EndpointUserConnections   = func(uID string) string { return EndpointUsers + uID + "/connections" }

	EndpointGuild                = func(gID string) string { return EndpointGuilds + gID }
	EndpointGuildInivtes         = func(gID string) string { return EndpointGuilds + gID + "/invites" }
	EndpointGuildChannels        = func(gID string) string { return EndpointGuilds + gID + "/channels" }
	EndpointGuildMembers         = func(gID string) string { return EndpointGuilds + gID + "/members" }
	EndpointGuildMember          = func(gID, uID string) string { return EndpointGuilds + gID + "/members/" + uID }
	EndpointGuildMemberRole      = func(gID, uID, rID string) string { return EndpointGuilds + gID + "/members/" + uID + "/roles/" + rID }
	EndpointGuildBans            = func(gID string) string { return EndpointGuilds + gID + "/bans" }
	EndpointGuildBan             = func(gID, uID string) string { return EndpointGuilds + gID + "/bans/" + uID }
	EndpointGuildIntegrations    = func(gID string) string { return EndpointGuilds + gID + "/integrations" }
	EndpointGuildIntegration     = func(gID, iID string) string { return EndpointGuilds + gID + "/integrations/" + iID }
	EndpointGuildIntegrationSync = func(gID, iID string) string { return EndpointGuilds + gID + "/integrations/" + iID + "/sync" }
	EndpointGuildRoles           = func(gID string) string { return EndpointGuilds + gID + "/roles" }
	EndpointGuildRole            = func(gID, rID string) string { return EndpointGuilds + gID + "/roles/" + rID }
	EndpointGuildInvites         = func(gID string) string { return EndpointGuilds + gID + "/invites" }
	EndpointGuildEmbed           = func(gID string) string { return EndpointGuilds + gID + "/embed" }
	EndpointGuildPrune           = func(gID string) string { return EndpointGuilds + gID + "/prune" }
	EndpointGuildIcon            = func(gID, hash string) string { return EndpointGuilds + gID + "/icons/" + hash + ".jpg" }
	EndpointGuildSplash          = func(gID, hash string) string { return EndpointGuilds + gID + "/splashes/" + hash + ".jpg" }
	EndpointGuildWebhooks        = func(gID string) string { return EndpointGuilds + gID + "/webhooks" }

	EndpointChannel                   = func(cID string) string { return EndpointChannels + cID }
	EndpointChannelPermissions        = func(cID string) string { return EndpointChannels + cID + "/permissions" }
	EndpointChannelPermission         = func(cID, tID string) string { return EndpointChannels + cID + "/permissions/" + tID }
	EndpointChannelInvites            = func(cID string) string { return EndpointChannels + cID + "/invites" }
	EndpointChannelTyping             = func(cID string) string { return EndpointChannels + cID + "/typing" }
	EndpointChannelMessages           = func(cID string) string { return EndpointChannels + cID + "/messages" }
	EndpointChannelMessage            = func(cID, mID string) string { return EndpointChannels + cID + "/messages/" + mID }
	EndpointChannelMessageAck         = func(cID, mID string) string { return EndpointChannels + cID + "/messages/" + mID + "/ack" }
	EndpointChannelMessagesBulkDelete = func(cID string) string { return EndpointChannel(cID) + "/messages/bulk_delete" }
	EndpointChannelMessagesPins       = func(cID string) string { return EndpointChannel(cID) + "/pins" }
	EndpointChannelMessagePin         = func(cID, mID string) string { return EndpointChannel(cID) + "/pins/" + mID }

	EndpointChannelWebhooks = func(cID string) string { return EndpointChannel(cID) + "/webhooks" }
	EndpointWebhook         = func(wID string) string { return EndpointWebhooks + wID }
	EndpointWebhookToken    = func(wID, token string) string { return EndpointWebhooks + wID + "/" + token }

	EndpointMessageReactions = func(cID, mID, eID string) string {
		return EndpointChannelMessage(cID, mID) + "/reactions/" + eID
	}
	EndpointMessageReaction = func(cID, mID, eID, uID string) string {
		return EndpointMessageReactions(cID, mID, eID) + "/" + uID
	}

	EndpointRelationships       = func() string { return EndpointUsers + "@me" + "/relationships" }
	EndpointRelationship        = func(uID string) string { return EndpointRelationships() + "/" + uID }
	EndpointRelationshipsMutual = func(uID string) string { return EndpointUsers + uID + "/relationships" }

	EndpointInvite = func(iID string) string { return EndpointAPI + "invite/" + iID }

	EndpointIntegrationsJoin = func(iID string) string { return EndpointAPI + "integrations/" + iID + "/join" }

	EndpointEmoji = func(eID string) string { return EndpointAPI + "emojis/" + eID + ".png" }

	EndpointOauth2          = EndpointAPI + "oauth2/"
	EndpointApplications    = EndpointOauth2 + "applications"
	EndpointApplication     = func(aID string) string { return EndpointApplications + "/" + aID }
	EndpointApplicationsBot = func(aID string) string { return EndpointApplications + "/" + aID + "/bot" }
)
