# Reference
## Chat
<details><summary><code>client.Chat.List() -> *roamhq.ListChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List accessible chats — DMs, multi-DMs, group chats, all-hands "team
Roam" groups, and meeting chats.

**Personal access tokens** are backed by the user's inbox: chats are
ordered by most recent activity and include `lastMessageTime`,
`isUnread`, `preview`, `isMuted`, and `isPinned`. Bot threads (where
the user has unread replies) are returned as separate rows keyed by
`threadTimestamp`.

**Organization tokens** receive the chats the bot has access to,
ordered by chat creation time. Inbox-derived fields
(`lastMessageTime`, `isUnread`, `preview`, `isMuted`, `isPinned`) are
not populated, since bot addresses do not accumulate inbox state for
normal messages — those are delivered via webhooks.

Timestamps are returned in the caller's timezone (see
[Timezone handling](https://developer.ro.am/docs/guides/migration-v0-to-v1#timezone-handling)).

**Required scope:** `chat:read`

Pass `expand=addresses` to include an address sidecar for chat participants
and preview senders. See [Identity & Principals](https://developer.ro.am/docs/guides/identity-and-principals).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListChatRequest{}
client.Chat.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**limit:** `*int` — Number of chats to return per response. Default 10, max 50.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>

<dl>
<dd>

**expand:** `*string` 

Comma-separated fields to expand. Supported: `addresses` — include an
`addresses` map resolving chat participants and preview sender IDs.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.Post(request) -> *roamhq.PostChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Send a message to a chat. Messages can be plain markdown text, rich [Block Kit](https://developer.ro.am/docs/guides/block-kit) layouts, or polls.

**Destination (ONE of the following is required):**
- `chatId` - Post to an existing chat by its ID
- `groupId` - Post to a group chat
- `userIds` - Post to a DM or Multi-DM with the specified users

You must specify exactly one destination. Specifying multiple destinations (e.g., both `chatId` and `groupId`) will return a 400 error.

Mentions use Slack's token syntax with Slack's semantics: `<@ID>` mentions a principal (a user or bot, e.g. `<@7861a4c6-765a-495d-898d-fae3d8fbba2d>` — resolvable via [`user.info`](https://developer.ro.am/docs/api/user-info)), `<!subteam^ID>` mentions a group or channel, notifying its members (resolvable via [`group.info`](https://developer.ro.am/docs/api/group-info)), and `<!channel>` notifies everyone in the chat.
When rendered in the client, the tag will automatically be replaced with the human-readable display name (or "everyone" for `<!channel>`).
On write, either token form is accepted for any mentionable ID; the legacy `<@all>` broadcast alias is accepted; and a Slack-style `|label` suffix (e.g. `<@7861a4c6-…|Rob>`, `<!subteam^59c1a4d2-…|@eng>`) is accepted and ignored — the mentioned entity's live display name is always used. Write-side acceptance is identical on every [API version](https://developer.ro.am/docs/guides/api-versioning). Messages read back always carry bare canonical tokens, and `<!channel>` for the broadcast — on API versions from `2026-08-07`; clients pinned to older versions read the older grammar (`<@ID>` for every mention, `<@all>`). Slack forms Roam does not implement are reserved and stay literal text: `<#ID>` channel links, `<!here>`, and `<!everyone>`.

**Custom sender (optional):** see the [Sender Profiles guide](https://developer.ro.am/docs/guides/sender-profiles).
- `sender.name` / `sender.imageUrl` are per-message display overrides, stored on the message itself.
- `sender.id` authors the message as a configured bot persona (Roam Administration > Developer > edit your app > Add Bot Persona). Ids that don't match a configured persona are accepted and ignored — the message is authored by the app's root identity. Sending never creates or renames personas.
- **Personal access tokens**: Reject the `sender` field with 400. PATs always post as their personal bot.

**Access:** Organization tokens can post to chats the bot is a member of,
and to **public groups** in the workspace without joining. Personal tokens
can post only where the owner is a member (`403` `not_in_chat` for an
unjoined public group). Full membership matrix:
[Chat](https://developer.ro.am/docs/guides/chat).

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.PostChatRequest{
    ChatID: roamhq.String(
        "757dfe66-37b4-4772-baa5-8c86ec68c176",
    ),
    Text: roamhq.String(
        "Hello from the **API**",
    ),
}
client.Chat.Post(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `*string` — Post to an existing chat by ID (mutually exclusive with groupId/userIds)
    
</dd>
</dl>

<dl>
<dd>

**groupID:** `*string` — Post to a group channel (mutually exclusive with chatId/userIds)
    
</dd>
</dl>

<dl>
<dd>

**userIDs:** `[]string` — Post to a DM or Multi-DM with these users (mutually exclusive with chatId/groupId)
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` 

Reply to a specific thread by providing the thread's timestamp.
If the timestamp doesn't correspond to an existing message, a 400 error is returned.
Mutually exclusive with `threadKey`.
    
</dd>
</dl>

<dl>
<dd>

**threadKey:** `*string` 

A stable external identifier used to group related messages into a thread.
On the first use of a given `threadKey`, a new message is posted and the resulting
thread timestamp is stored. Subsequent messages with the same `threadKey` are
automatically threaded under the original message.

This is useful for external integrations (e.g. PagerDuty, Grafana, Sentry) that
want to thread related messages using their own identifiers (such as `dedup_key`,
`fingerprint`, or `group_id`) without tracking Roam's internal thread timestamps.

Mutually exclusive with `threadTimestamp`. When `threadKey` is provided, the
response is always synchronous (equivalent to `sync: true`).
    
</dd>
</dl>

<dl>
<dd>

**replyTimestamp:** `*int` 

Reply directly to a specific message by its timestamp. Unlike
`threadTimestamp` (which threads a reply under a parent message in a
group), `replyTimestamp` is a direct reply used in DMs — which have no
threads — and within an existing channel thread. Text messages only:
not supported together with `blocks` or `poll`.
    
</dd>
</dl>

<dl>
<dd>

**text:** `*string` — Message text in GitHub-flavored markdown
    
</dd>
</dl>

<dl>
<dd>

**markdown:** `*bool` — Text is markdown by default. If set to false, markdown interpretation will be disabled.
    
</dd>
</dl>

<dl>
<dd>

**items:** `[]string` — Array of Item IDs to attach to this message.
    
</dd>
</dl>

<dl>
<dd>

**assetIDs:** `[]string` 

Array of asset IDs from [`/asset.create`](https://developer.ro.am/docs/api/asset-create)
to attach to this message. Each asset must be owned by your app
and fully uploaded (processed and ready). Combines with
`text`/`items`; not with `blocks` or `poll`.
    
</dd>
</dl>

<dl>
<dd>

**blocks:** `[]*roamhq.PostChatRequestBlocksItem` 

Array of [Block Kit](https://developer.ro.am/docs/guides/block-kit) block objects for rich message formatting.
Cannot be combined with `text` or `items`. Maximum 10 blocks, 8,000 bytes total payload.
    
</dd>
</dl>

<dl>
<dd>

**color:** `*string` 

Colored vertical strip on the side of the message. Only used with `blocks`.
Named values: `good` (green), `warning` (yellow), `danger` (red), or a hex color like `#5B3FD9`.
    
</dd>
</dl>

<dl>
<dd>

**poll:** `*roamhq.PostChatRequestPoll` — Create a poll message. Mutually exclusive with `text`, `items`, and `blocks`.
    
</dd>
</dl>

<dl>
<dd>

**sender:** `*roamhq.Sender` 
    
</dd>
</dl>

<dl>
<dd>

**sync:** `*bool` — If set, the post will be performed synchronously and its timestamp returned. Incompatible with `sendAt`.
    
</dd>
</dl>

<dl>
<dd>

**sendAt:** `*time.Time` 

Schedule the message for later delivery (RFC 3339). Requirements:
- Must be in the **future** and within **30 days**
- Must fall on a **15-minute UTC boundary** (`:00`, `:15`, `:30`, or `:45`; seconds and sub-seconds zero)
- Incompatible with `sync`, `poll`, `threadKey`, and `replyTimestamp`

When `sendAt` is set, the response is `{chatId, scheduledMessageId, sendAt}`
instead of an immediate message `timestamp`.

Scheduled messages can be listed via
[`/chat.scheduled.list`](https://developer.ro.am/docs/api/chat-scheduled-list) and canceled via
[`/chat.scheduled.cancel`](https://developer.ro.am/docs/api/chat-scheduled-cancel) until they send.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.PostEphemeral(request) -> *roamhq.PostEphemeralChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Post an **ephemeral message** — visible to a single member of a chat, with an
"Only you can see this" header — without posting anything the other members can
see. This is the standard way for a bot to respond privately in a shared
channel (the Roam equivalent of Slack's `chat.postEphemeral`).

The target `userId` must be a member of the chat (for channels: a member of the
backing group), otherwise the request fails with `user_not_in_chat`.

`text` is always rendered as GitHub-flavored markdown. Mention markup
(`<@USER_ID>`) is **not** supported in ephemeral messages. Block Kit `blocks`
are not currently supported.

**Delivery semantics — read before using:**
- **Desktop and web only.** Mobile clients do not display ephemeral messages,
  and no mobile push notification is sent. A recipient who only uses Roam on
  mobile will never see the message.
- **Best-effort, at-most-once.** The message is delivered in real time to the
  recipient's connected clients, and to recently-active offline clients when
  they reconnect. A recipient who has been offline for several days (or has
  never signed in on that device) silently misses it. There are no retries
  and no delivery receipt.
- **Transient.** The message is never stored server-side. It disappears when
  the recipient restarts their app, and it never appears in
  [`/chat.history`](https://developer.ro.am/docs/api/chat-history) or [`/chat.search`](https://developer.ro.am/docs/api/chat-search).
- **Not addressable.** It cannot be edited or deleted:
  [`/chat.update`](https://developer.ro.am/docs/api/chat-update) and [`/chat.delete`](https://developer.ro.am/docs/api/chat-delete)
  against its `(chatId, timestamp)` return `message_not_found`.
- **No webhooks.** Posting an ephemeral message never triggers a
  [`chat.message`](https://developer.ro.am/docs/webhooks/chat-message) event, so it cannot leak to
  org-wide webhook consumers.

Do not use ephemeral messages for anything the recipient must durably receive —
use a DM ([`/chat.post`](https://developer.ro.am/docs/api/chat-post) with `userIds`) for that.

**Custom sender (optional):** same semantics as [`/chat.post`](https://developer.ro.am/docs/api/chat-post) —
`sender.name` / `sender.imageUrl` apply a per-message display override, and
`sender.id` authors the message as a configured bot persona (unknown ids
are accepted and ignored). Personal access tokens reject the `sender`
field. See the [Sender Profiles guide](https://developer.ro.am/docs/guides/sender-profiles).

**Required scope:** `chat:send_message` or `chat:write`

**Access:** Organization and Personal. The organization bot or
personal-token **owner** must be a member of the chat (`403` `not_in_chat`
otherwise) — unlike [`/chat.post`](https://developer.ro.am/docs/api/chat-post), there is no
public-group carveout. Personal tokens send as the user's personal bot
and reject the `sender` field.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.PostEphemeralChatRequest{
    ChatID: "295155ae-7df5-4ed5-9ebc-89a170559c81",
    UserID: "7861a4c6-765a-495d-898d-fae3d8fbba2d",
    Text: "Only *you* can see this: your deploy token expires in 3 days.",
}
client.Chat.PostEphemeral(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `string` — The chat to post into. Use [`/chat.list`](https://developer.ro.am/docs/api/chat-list) or a `chat.message` webhook payload to obtain chat IDs.
    
</dd>
</dl>

<dl>
<dd>

**userID:** `string` — The user who should see the message. Must be a member of the chat.
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` 

Show the ephemeral message inside an existing thread instead of the
main channel view. Channels only — returns 400 in DMs and Multi-DMs.
The value is not validated against an existing thread: pass a real
thread's timestamp, or the message is keyed under a thread view the
recipient can never open and is effectively never seen.
    
</dd>
</dl>

<dl>
<dd>

**text:** `string` 

Message text in GitHub-flavored markdown (always rendered as
markdown; there is no plain-text mode). Maximum 8,000 bytes.
Mention markup is not supported.
    
</dd>
</dl>

<dl>
<dd>

**sender:** `*roamhq.Sender` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.ListScheduled() -> *roamhq.ListScheduledChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists pending messages scheduled via [`/chat.post`](https://developer.ro.am/docs/api/chat-post)'s `sendAt`
that have not been sent yet. Results are ordered ascending by `sendAt` (soonest
first). Sent and canceled messages are not returned.

Only messages scheduled by the calling credential's bot identity are listed:
organization tokens of the same app share the app's bot identity (and therefore
see each other's scheduled messages), while personal access tokens have a
per-person bot identity and see only their own.

**Access:** Organization and Personal.

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListScheduledChatRequest{}
client.Chat.ListScheduled(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `*string` — Only return messages scheduled for this chat.
    
</dd>
</dl>

<dl>
<dd>

**after:** `*time.Time` 

Only return messages scheduled to send after this datetime
(YYYY-MM-DD or RFC-3339). Exclusive.
    
</dd>
</dl>

<dl>
<dd>

**before:** `*time.Time` 

Only return messages scheduled to send before this datetime
(YYYY-MM-DD or RFC-3339). Exclusive.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — The number of scheduled messages to return per response. Default is 10.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.CancelScheduled(request) -> *roamhq.CancelScheduledChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Cancels a pending message scheduled via [`/chat.post`](https://developer.ro.am/docs/api/chat-post)'s
`sendAt`, so it will never be delivered. Pending scheduled messages can be
discovered with [`/chat.scheduled.list`](https://developer.ro.am/docs/api/chat-scheduled-list).

Only the credential's bot identity that scheduled the message may cancel it. A
`scheduledMessageId` scheduled by a different identity — or one that never
existed — returns `scheduled_message_not_found`; the endpoint does not reveal
whether such an id exists. Canceling a message that has already been sent
returns `scheduled_message_already_sent`.

Cancellation is best-effort once the scheduled send time arrives: delivery of a
due message begins in the seconds after its `sendAt` boundary, and a cancel
issued inside that window may return success while the message is still
delivered. Cancel ahead of the scheduled time to be safe.

**Access:** Organization and Personal.

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.CancelScheduledChatRequest{
    ScheduledMessageID: "0197f9f0-5cc1-7d07-8a12-9e65a8a0c1b9",
}
client.Chat.CancelScheduled(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**scheduledMessageID:** `string` — The id returned by `/chat.post` when the message was scheduled.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.StartStream(request) -> *roamhq.StartStreamChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Open a streaming message and post its first content. Streaming lets a bot
deliver a message incrementally — recipients see the text fill in live (with
a "typing…" indicator) instead of waiting for the full response. This is
useful for AI agents that produce text token-by-token.

A stream has three steps, each its own request:

1. **[`/chat.startStream`](https://developer.ro.am/docs/api/chat-start-stream)** — open the stream and pick the destination. Returns a `streamId`.
2. **[`/chat.appendStream`](https://developer.ro.am/docs/api/chat-append-stream)** — append chunks of text (call as many times as needed).
3. **[`/chat.stopStream`](https://developer.ro.am/docs/api/chat-stop-stream)** — finalize the stream into a single persisted message.

Pass the `streamId` returned here to every subsequent `appendStream` and
`stopStream`. The sender, destination, and thread are fixed for the lifetime
of the stream.

**Custom sender (optional):** same semantics as
[`/chat.post`](https://developer.ro.am/docs/api/chat-post) — `sender.name` / `sender.imageUrl`
apply a per-message display override to the finalized message, and
`sender.id` authors the stream as a configured bot persona (unknown ids
are accepted and ignored). The typing indicator shown while streaming uses
the override name when given, otherwise the persona's or app's configured
name. See the [Sender Profiles guide](https://developer.ro.am/docs/guides/sender-profiles).

**Access:** Organization and Personal. Organization tokens follow the
same public-group carveout as [`/chat.post`](https://developer.ro.am/docs/api/chat-post): the
bot may stream into a public group in its roam without joining.
Personal tokens can stream only where the owner is a member
(`403` `not_in_chat` for an unjoined public group) and reject the
`sender` field.

**Required scope:** `chat:send_message` or `chat:write`

## Destination

Provide exactly one of `chatId`, `groupId`, or `userIds`. If `text` is empty,
the destination is recorded but message creation is deferred until the first
non-empty `appendStream` or the `stopStream` call.

## Thinking streams

Set `kind` to `thinking` to finalize the message as a thought-bubble; clients
show a "thinking…" indicator instead of "typing…". The default `kind` is `text`.

## Limits

- Up to **10 concurrent streams per API client**.
- Only **one active stream per chat** at a time.
- Accumulated text may not exceed the regular message size limit.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.StartStreamChatRequest{
    GroupID: roamhq.String(
        "88bebce7-6cbb-4666-96f9-5c02d73e6661",
    ),
    Text: roamhq.String(
        "Let me look into that...",
    ),
}
client.Chat.StartStream(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `*string` — Stream into an existing chat by ID (mutually exclusive with groupId/userIds).
    
</dd>
</dl>

<dl>
<dd>

**groupID:** `*string` — Stream into a group chat (mutually exclusive with chatId/userIds).
    
</dd>
</dl>

<dl>
<dd>

**userIDs:** `[]string` — Stream into a DM or Multi-DM with these users (mutually exclusive with chatId/groupId).
    
</dd>
</dl>

<dl>
<dd>

**kind:** `*roamhq.StartStreamChatRequestKind` — Stream kind. `thinking` finalizes as a thought-bubble message.
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` — Optional thread to reply within.
    
</dd>
</dl>

<dl>
<dd>

**text:** `*string` — Optional initial text. May be empty to defer destination resolution until the first append/stop.
    
</dd>
</dl>

<dl>
<dd>

**sender:** `*roamhq.Sender` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.AppendStream(request) -> *roamhq.AppendStreamChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Append a chunk of text to an open stream (see
[`/chat.startStream`](https://developer.ro.am/docs/api/chat-start-stream)). Each chunk is broadcast
to recipients as a delta, so the message appears to fill in live. Call as
many times as needed before [`/chat.stopStream`](https://developer.ro.am/docs/api/chat-stop-stream).

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.AppendStreamChatRequest{
    StreamID: "018f5c8e-7d2a-7c4e-8f9a-1a2b3c4d5e6f",
    Text: " The answer is 42.",
}
client.Chat.AppendStream(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**streamID:** `string` — The stream ID returned by chat.startStream.
    
</dd>
</dl>

<dl>
<dd>

**text:** `string` — Text chunk to append. Required and non-empty.
    
</dd>
</dl>

<dl>
<dd>

**snapshot:** `*bool` 

If `true`, **replace** the accumulated text with `text` (and broadcast it
as a full snapshot) instead of appending. Useful when the client holds the
canonical current state — for example after rewriting prior output. The
message size limit is applied to the new `text` alone.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.StopStream(request) -> *roamhq.StopStreamChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Finalize an open stream (see [`/chat.startStream`](https://developer.ro.am/docs/api/chat-start-stream))
into a single persisted chat message and return its timestamp. Optionally
include trailing `text` to append before finalizing.

If the app never calls `stopStream` but has already streamed some text, the
server finalizes the buffered text into a message automatically.

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.StopStreamChatRequest{
    StreamID: "018f5c8e-7d2a-7c4e-8f9a-1a2b3c4d5e6f",
    Text: roamhq.String(
        " Hope that helps!",
    ),
}
client.Chat.StopStream(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**streamID:** `string` — The stream ID returned by chat.startStream.
    
</dd>
</dl>

<dl>
<dd>

**text:** `*string` — Optional trailing text appended before the message is finalized.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.Update(request) -> *roamhq.UpdateChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Edit a previously posted bot message. The updated message can contain plain markdown text or rich [Block Kit](https://developer.ro.am/docs/guides/block-kit) layouts.

The bot must own the message being updated (matched by address ID). Personal access tokens always send as their bot persona and may only edit messages that personal bot posted.

**Access:** Organization and Personal.

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.UpdateChatRequest{
    ChatID: "757dfe66-37b4-4772-baa5-8c86ec68c176",
    Timestamp: 1765602474760032,
    Text: roamhq.String(
        "Updated message content with **bold text**",
    ),
}
client.Chat.Update(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `string` — ID of the chat containing the message.
    
</dd>
</dl>

<dl>
<dd>

**timestamp:** `int` — Timestamp of the message to update.
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` — Thread timestamp, if the message is in a thread.
    
</dd>
</dl>

<dl>
<dd>

**text:** `*string` 

Updated markdown-formatted text content. Required unless `blocks` is provided.
Cannot be combined with `blocks`.
    
</dd>
</dl>

<dl>
<dd>

**markdown:** `*bool` — Text is markdown by default. If this is set to false, markdown interpretation will be disabled.
    
</dd>
</dl>

<dl>
<dd>

**items:** `[]string` — Array of Item IDs to attach to this message. Cannot be combined with `blocks`.
    
</dd>
</dl>

<dl>
<dd>

**assetIDs:** `[]string` 

Array of asset IDs from [`/asset.create`](https://developer.ro.am/docs/api/asset-create)
to attach to this message. Each asset must be owned by your app
and fully uploaded (processed and ready). Cannot be combined with `blocks`.
    
</dd>
</dl>

<dl>
<dd>

**blocks:** `[]*roamhq.UpdateChatRequestBlocksItem` 

Array of [Block Kit](https://developer.ro.am/docs/guides/block-kit) block objects for rich message formatting.
Cannot be combined with `text` or `items`. Maximum 10 blocks, 8,000 bytes total payload.
    
</dd>
</dl>

<dl>
<dd>

**color:** `*string` 

Colored vertical strip on the side of the message. Only used with `blocks`.
Named values: `good` (green), `warning` (yellow), `danger` (red), or a hex color like `#5B3FD9`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.Delete(request) -> *roamhq.DeleteChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Delete a previously posted bot message. The bot must own the message being deleted (matched by address ID). Personal access tokens always send as their bot persona and may only delete messages that personal bot posted.

Deleting an already-deleted message is idempotent and returns success.

**Access:** Organization and Personal.

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.DeleteChatRequest{
    ChatID: "757dfe66-37b4-4772-baa5-8c86ec68c176",
    Timestamp: 1765602474760032,
}
client.Chat.Delete(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `string` — ID of the chat containing the message.
    
</dd>
</dl>

<dl>
<dd>

**timestamp:** `int` — Timestamp of the message to delete.
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` — Thread timestamp, if the message is in a thread.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.Typing(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Notify other chat participants that you are working on a response.
If they have the chat open, they will see "(Bot name) is typing...".

The indicator lasts **6 seconds**. Re-send every **5 seconds** to keep
it visible while you work. Longer gaps will let it expire between pings.

**Destination options (mutually exclusive):**
- `chatId` - Send to an existing chat by its ID
- `groupId` - Send to a group channel
- `userIds` - Send to a DM or Multi-DM with the specified users

**Custom sender (optional):** pass `sender.id` to show the indicator as a
[configured bot persona](https://developer.ro.am/docs/guides/sender-profiles) — the persona's
configured name and avatar are used. Only `id` is accepted; `name` and
`imageUrl` are rejected on this endpoint. Selection is lookup-only: an id
that doesn't match a configured persona is accepted and ignored, and the
indicator shows the app's own identity (same for an omitted, empty, or `_`
id). Personal access tokens reject `sender` entirely.

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.TypingChatRequest{
    ChatID: roamhq.String(
        "295155ae-7df5-4ed5-9ebc-89a170559c81",
    ),
}
client.Chat.Typing(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `*string` — Send to an existing chat by ID (mutually exclusive with groupId/userIds)
    
</dd>
</dl>

<dl>
<dd>

**groupID:** `*string` — Send to a group channel (mutually exclusive with chatId/userIds)
    
</dd>
</dl>

<dl>
<dd>

**userIDs:** `[]string` — Send to a DM or Multi-DM with these users (mutually exclusive with chatId/groupId)
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` — Timestamp of the message being replied to.
    
</dd>
</dl>

<dl>
<dd>

**sender:** `*roamhq.TypingChatRequestSender` 

Optional configured bot persona to show the indicator as. Only
`id` is accepted — `name` and `imageUrl` are rejected on this
endpoint. Personal access tokens reject this field entirely.
See the [Sender Profiles guide](https://developer.ro.am/docs/guides/sender-profiles).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.History() -> *roamhq.HistoryChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List messages in a chat, filtered by date range (after/before).

Messages with `contentType` of `text`, `voice`, or `poll` are returned. System messages and other content types are excluded.

**Specify ONE of the following:**
- `chatId` - Fetch from an existing chat by its ID
- `groupId` - Fetch from a group chat
- `userIds` - Fetch from a DM or Multi-DM with the specified users

You must specify exactly one destination. Specifying multiple (e.g., both `chatId` and `groupId`) will return a 400 error.

The ordering of results depends on the filter specified:

- When no parameters are provided, the most recent messages are returned,
  sorted in reverse chronological order. This is equivalent to specifying `before`
  as NOW and leaving `after` unspecified.

- If `after` is specified, the results are sorted in forward chronological order.

Either dates or datetimes may be specified. Date-only inputs (`YYYY-MM-DD`)
are interpreted in the caller's timezone (see
[Timezone handling](https://developer.ro.am/docs/guides/migration-v0-to-v1#timezone-handling)).

**Access:** Organization tokens need to be a **member** of the chat
(`403` `not_in_chat` otherwise). Personal tokens can read any chat the
owner can, including public groups in their roam they have not joined.
Full membership matrix: [Chat](https://developer.ro.am/docs/guides/chat).

**Required scope:** `chat:history`

Every returned sender includes `userId` plus `userType`. The ID resolves
through [`user.info`](https://developer.ro.am/docs/api/user-info) with the same credentials.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.HistoryChatRequest{}
client.Chat.History(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `*string` — The chat ID to fetch messages from. Either chatId, groupId, or userIds must be specified.
    
</dd>
</dl>

<dl>
<dd>

**groupID:** `*string` — Group chat ID to fetch messages from. Either chatId, groupId, or userIds must be specified.
    
</dd>
</dl>

<dl>
<dd>

**userIDs:** `*string` — User IDs to fetch DM/Multi-DM messages with. Either chatId, groupId, or userIds must be specified.
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*float64` — Read replies of the message with this timestamp. Specified in microseconds.
    
</dd>
</dl>

<dl>
<dd>

**after:** `*string` 

The datetime to begin listing messages (YYYY-MM-DD or RFC-3339).
Date-only values are interpreted in the caller's timezone.
Sub-millisecond precision on datetimes is truncated. Defaults to
"no filter".
    
</dd>
</dl>

<dl>
<dd>

**before:** `*string` 

The datetime until which to list messages (YYYY-MM-DD or RFC-3339).
Date-only values are interpreted in the caller's timezone.
Sub-millisecond precision on datetimes is truncated. Defaults to
"now".
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Number of messages to return (default 10, max 200).
    
</dd>
</dl>

<dl>
<dd>

**expand:** `*string` 

Comma-separated fields to expand. Supported: `addresses` — include an
`addresses` map resolving the sender (`userId`) and mentioned IDs on
each message to their display info.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.Search(request) -> *roamhq.SearchChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Full-text search over the caller's accessible messages. Returns
full-fidelity messages — text, items, voice, polls, blocks, and
mentions — hydrated through the same pipeline as
[`/chat.history`](https://developer.ro.am/docs/api/chat-history).

All fields are optional. With no parameters, the most recent messages
across all chat types (DMs, multi-DMs, group chats) are returned in
reverse chronological order.

**Sort:** When omitted and `query` is empty, results are sorted
chronologically (newest first), since relevance scoring is meaningless
without search terms. Pass `sort: recent` to force chronological order
even with a text query.

**Date filters:** `before` and `after` accept `YYYY-MM-DD`. Dates are
interpreted in the caller's timezone (see
[Timezone handling](https://developer.ro.am/docs/guides/migration-v0-to-v1#timezone-handling)).

**Access:** Organization and Personal.

- **Personal tokens** search chats the owner can read, including public
  groups in their roam they have not joined.
- **Organization tokens** search chats the bot is a **member** of,
  plus unjoined **public** groups in the bot's roam (Slack
  `search:read.public`). Private groups the bot is not in are excluded.
  [`/chat.history`](https://developer.ro.am/docs/api/chat-history) stays membership-only.

Full membership matrix: [Chat](https://developer.ro.am/docs/guides/chat).

**Required scope:** `chat:history`

Every returned sender includes `userId` plus `userType`. The ID resolves
through [`user.info`](https://developer.ro.am/docs/api/user-info) with the same credentials.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.SearchChatRequest{
    After: roamhq.String(
        "2026-04-13",
    ),
    Limit: roamhq.Int(
        20,
    ),
}
client.Chat.Search(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**query:** `*string` — Free-text search query. Empty matches all messages.
    
</dd>
</dl>

<dl>
<dd>

**in:** `[]string` — Group names to search within.
    
</dd>
</dl>

<dl>
<dd>

**from:** `[]string` — Filter to messages sent by these email addresses.
    
</dd>
</dl>

<dl>
<dd>

**with:** `[]string` — Filter to chats including these email addresses.
    
</dd>
</dl>

<dl>
<dd>

**before:** `*string` — Only include messages before this date (`YYYY-MM-DD`, caller's timezone).
    
</dd>
</dl>

<dl>
<dd>

**after:** `*string` — Only include messages on or after this date (`YYYY-MM-DD`, caller's timezone).
    
</dd>
</dl>

<dl>
<dd>

**has:** `[]*roamhq.SearchChatRequestHasItem` — Restrict to messages that contain a mention or an item.
    
</dd>
</dl>

<dl>
<dd>

**chatTypes:** `[]*roamhq.SearchChatRequestChatTypesItem` 

Restrict to specific chat types. Defaults to all types
(channels, all-hands "team Roam" groups, and DMs).
    
</dd>
</dl>

<dl>
<dd>

**excludeChatIDs:** `[]string` — Chat IDs to exclude from results.
    
</dd>
</dl>

<dl>
<dd>

**excludeUserIDs:** `[]string` — Sender user IDs to exclude from results.
    
</dd>
</dl>

<dl>
<dd>

**sort:** `*roamhq.SearchChatRequestSort` 

`relevant` (default) ranks by relevance to `query`; `recent`
sorts newest first. With an empty `query`, results are
sorted chronologically regardless.
    
</dd>
</dl>

<dl>
<dd>

**expand:** `*string` 

Comma-separated fields to expand. Supported: `addresses` —
include an `addresses` map resolving the sender (`userId`) and
mentioned IDs on each message to their display info.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Number of messages per page (max 200).
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.ResolveLink(request) -> *roamhq.ResolveLinkChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Parse a Roam chat deep link (e.g. `https://ro.am/r/#/d/...`) and return the
referenced message.

When the caller has access to the referenced chat, the full message is
returned and `readable` is `true`. The `message` object is the same
shape as a `chat.history`/`chat.search` message — same fields, same
mention rendering. When the caller lacks access, the response still
includes the message key (`chatId`, `timestamp`, and `threadTimestamp`
if applicable) with `readable: false` and no message content — suitable
for rendering a reference without leaking content.

Use [`/chat.link.create`](https://developer.ro.am/docs/api/chat-link-create) for the reverse
operation — minting a shareable Roam link from a message the caller can
already read.

**Access:** Organization and Personal.

**Required scope:** `chat:history`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ResolveLinkChatRequest{
    Link: "https://ro.am/r/#/d/abc123xyz/c/757dfe66-37b4-4772-baa5-8c86ec68c176?ts=1765602474760032",
}
client.Chat.ResolveLink(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**link:** `string` — A Roam chat deep link URL that contains a message reference.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.CreateLink(request) -> *roamhq.CreateLinkChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Create a shareable Roam link to a specific chat message. Opening the link
in Roam navigates to that message in its chat.

Identify the chat with exactly one of `chatId`, `groupId`, or `userIds`,
and the message by its `timestamp` (Unix microseconds), as returned by
[`/chat.history`](https://developer.ro.am/docs/api/chat-history), [`/chat.post`](https://developer.ro.am/docs/api/chat-post),
or webhook message events. For a thread reply, also pass the thread root's
timestamp as `threadTimestamp` — without it the reply will not be found.

The message must exist and be readable by the caller; otherwise no link is
returned (`404` if the message does not exist, `403` if the caller is not a
member of the chat). The link itself does not grant access: recipients can
only open it if they are members of the chat.

Use [`/chat.link.resolve`](https://developer.ro.am/docs/api/chat-link-resolve) for the reverse
operation — turning a Roam chat link back into the referenced message.

**Access:** Organization and Personal. In Personal mode, only chats the
authenticated user can access are allowed.

**Required scope:** `chat:history`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.CreateLinkChatRequest{
    ChatID: roamhq.String(
        "295155ae-7df5-4ed5-9ebc-89a170559c81",
    ),
    Timestamp: int64(1765602474760032),
}
client.Chat.CreateLink(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `*string` — ID of the chat containing the message. Exactly one of `chatId`, `groupId`, or `userIds` is required.
    
</dd>
</dl>

<dl>
<dd>

**groupID:** `*string` — ID of a group whose channel chat contains the message.
    
</dd>
</dl>

<dl>
<dd>

**userIDs:** `[]string` — User ID(s) identifying the DM or group DM containing the message.
    
</dd>
</dl>

<dl>
<dd>

**timestamp:** `int64` — The message's timestamp in Unix microseconds.
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int64` — The thread root's timestamp in Unix microseconds. Required when the message is a thread reply.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Chat.Unfurl(request) -> *roamhq.UnfurlChatResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Attach app-provided preview cards to links in an existing text message.
Every map key must be an exact URL currently present in the message and
must match one of the app's registered unfurl domains. Validation is
atomic: if any entry is invalid, no previews are changed.

App previews replace Roam-generated previews for the same exact URL while
preserving unrelated previews. The server does not fetch any URL supplied
in this request.

**Access:** Organization only (API Key or OAuth). Register unfurl domains on
the API client first — see [Unfurling links](https://developer.ro.am/docs/guides/unfurling-links).
Personal Access Tokens cannot register domains or call this endpoint.

**Required scope:** `links:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.UnfurlChatRequest{
    ChatID: "8f3b9c2e-1a4d-4e7b-9c0a-2b6d1f5e3a7c",
    MessageTimestamp: 1748906400000000,
    Unfurls: map[string]*roamhq.UnfurlContent{
        "https://status.example.com/incidents/123": &roamhq.UnfurlContent{
            Title: "Incident 123",
            Description: roamhq.String(
                "Investigating elevated errors",
            ),
            SiteName: roamhq.String(
                "PagerDuty",
            ),
            Favicon: roamhq.String(
                "https://status.example.com/favicon.png",
            ),
            Image: &roamhq.UnfurlContentImage{
                URL: "https://status.example.com/incident.png",
                Type: roamhq.String(
                    "image/png",
                ),
                Width: roamhq.Int(
                    1200,
                ),
                Height: roamhq.Int(
                    630,
                ),
                Alt: roamhq.String(
                    "Incident status",
                ),
            },
        },
    },
}
client.Chat.Unfurl(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**messageTimestamp:** `int` — Timestamp of a top-level or threaded message in Unix microseconds.
    
</dd>
</dl>

<dl>
<dd>

**unfurls:** `map[string]*roamhq.UnfurlContent` — Preview content keyed by the exact URL from the message.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Reaction
<details><summary><code>client.Reaction.Add(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Add a reaction to a message in a chat.

To react to a thread reply, provide the `threadTimestamp` of the parent message
and the `timestamp` of the specific reply.

**Access:** The organization bot or personal-token **owner** must be a
member of the chat (`403` `not_in_chat` otherwise).

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.AddReactionRequest{
    ChatID: "7be17589-4b9a-4524-bddb-ce60abea08e6",
    Timestamp: 1755723832718034,
    Name: "thumbs_up",
}
client.Reaction.Add(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `string` — The chat containing the message to react to.
    
</dd>
</dl>

<dl>
<dd>

**timestamp:** `int` — Timestamp of the message to react to (Unix microseconds).
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` — Timestamp of the parent thread message (Unix microseconds), if reacting to a thread reply.
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — Name of the reaction to add (e.g. "thumbs_up", "heart", "100").
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Reaction.Remove(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove a reaction from a message in a chat.

Only reactions added by the authenticated app can be removed.

To remove a reaction from a thread reply, provide the `threadTimestamp` of the parent message
and the `timestamp` of the specific reply.

**Required scope:** `chat:send_message` or `chat:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.RemoveReactionRequest{
    ChatID: "7be17589-4b9a-4524-bddb-ce60abea08e6",
    Timestamp: 1755723832718034,
    Name: "thumbs_up",
}
client.Reaction.Remove(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `string` — The chat containing the message.
    
</dd>
</dl>

<dl>
<dd>

**timestamp:** `int` — Timestamp of the message (Unix microseconds).
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` — Timestamp of the parent thread message (Unix microseconds), if removing from a thread reply.
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — Name of the reaction to remove (e.g. "thumbs_up", "heart").
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Reaction.List() -> *roamhq.ListReactionResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List reactions on a specific message, grouped by emoji (Slack-style
`{name, count, users}`). Poll votes are returned separately in `pollVotes`
rather than folded into `reactions`.

`users` contains visible principal IDs only. Unknown or unauthorized actors
are omitted, and `count` is recomputed from the returned IDs. Hydrate them
with `user.list?ids`; these arrays do not carry inline type fields.

To list reactions on a thread reply, provide the `threadTimestamp` of the
parent message and the `timestamp` of the specific reply.

**Required scope:** `chat:history`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListReactionRequest{
    ChatID: "chatId",
    Timestamp: 1,
}
client.Reaction.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**chatID:** `string` — The chat containing the message.
    
</dd>
</dl>

<dl>
<dd>

**timestamp:** `int` — Timestamp of the message (Unix microseconds).
    
</dd>
</dl>

<dl>
<dd>

**threadTimestamp:** `*int` — Timestamp of the parent thread message (Unix microseconds), if listing reactions on a thread reply.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Asset
<details><summary><code>client.Asset.Create(request) -> *roamhq.CreateAssetResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Create a file asset and get back a self-describing instruction for
uploading its bytes — the JSON-friendly way to attach a file (image, PDF,
document, …) to a message, supply media for a story, or host an avatar
image. Unlike [`/item.upload`](https://developer.ro.am/docs/api/item-upload), which takes raw
bytes in the request body, every caller-visible step here is JSON in /
JSON out (so it can be driven from MCP and other tool-calling clients),
and the file bytes never pass through this API.

**Flow:**
1. `POST /asset.create` with the file `name` (include the extension, e.g.
   `photo.png`) and, if known, its `size` in bytes. For stories, also pass
   `purpose: "story"`. For avatars, pass `purpose: "avatar"` and `size`
   (max 10 MiB). The response
   is an upload instruction: `assetId`, `uploadUrl`, `uploadMethod`, and
   `uploadHeaders`. Avatar responses also include `imageUrl`.
2. Upload the raw bytes in a **single request**: use `uploadMethod` (a
   `POST`) against `uploadUrl`, send every header from `uploadHeaders`
   verbatim, and put the file in the request body. Send the headers exactly
   as given — they authorize the upload and select the single-request
   upload protocol; omitting any will cause the upload to fail.
3. Processing (thumbnails, previews, 512×512 WebP for avatars) happens
   automatically once the bytes land. There is no separate "complete" call.
4. Once the asset is ready, use it:
   - `purpose: "file"` (default) — attach via `assetIds` on
     [`/chat.post`](https://developer.ro.am/docs/api/chat-post) or
     [`/chat.update`](https://developer.ro.am/docs/api/chat-update)
   - `purpose: "story"` — post via [`/story.post`](https://developer.ro.am/docs/api/story-post)
   - `purpose: "avatar"` — pass `imageUrl` as `sender.imageUrl` on
     [`/chat.post`](https://developer.ro.am/docs/api/chat-post) (and related send endpoints), or
     as `hosts[].imageUrl` on
     [`/onair.event.create`](https://developer.ro.am/docs/onair-api/onair-event-create) /
     [`/onair.event.update`](https://developer.ro.am/docs/onair-api/onair-event-update)

A freshly-uploaded asset may take a few seconds to process (videos take
longer). Chat and story endpoints that consume the asset return a 400 with
a "still processing" message until processing completes. Avatar `imageUrl`
404s until the image is ready — wait a moment after the upload returns
before posting it.

The `uploadUrl` is short-lived; if it expires, call `asset.create` again for
a fresh instruction. Maximum file size is 5 GiB for `file` / `story`, and
10 MiB for `avatar`.

## Purposes

| Purpose | Use | Access |
|---------|-----|--------|
| `file` (default) | Chat message attachments | Organization and Personal |
| `story` | Story media (photo or video) | Personal only |
| `avatar` | `sender.imageUrl` and On-Air `hosts.imageUrl` | Organization and Personal |

Story assets are owned by the authenticated user (stories are posted as you,
not as a bot) and expire about 48 hours after creation. Because the media
must outlive the story's 24-hour lifetime, call
[`/story.post`](https://developer.ro.am/docs/api/story-post) within about 23 hours of creating the
asset; after that the asset is rejected and a new one must be created.

Avatar assets are public 512×512 WebP images. They do not expire. From
API version `2026-08-25`, `sender.imageUrl` and On-Air `hosts.imageUrl`
must be a Roam-hosted avatar URL (this `imageUrl`, or a legacy
`/card-images/` or `/photos/people/` URL). Third-party image URLs return
400. See [API Versioning](https://developer.ro.am/docs/guides/api-versioning) and
[Sender Profiles](https://developer.ro.am/docs/guides/sender-profiles).

**Access:** Organization and Personal. `purpose: "story"` is Personal only.

**Required scope:** `item:write` for `purpose: "file"`; `chat:send_message`
or `chat:write` for `purpose: "story"`; any of `item:write`,
`chat:send_message`, `chat:write`, or `onair:write` for `purpose: "avatar"`.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.CreateAssetRequest{
    Name: "quarterly-report.pdf",
    Size: roamhq.Int(
        248173,
    ),
}
client.Asset.Create(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**name:** `string` — File name, including its extension (e.g. `report.pdf`). Processing determines the media type from the extension.
    
</dd>
</dl>

<dl>
<dd>

**size:** `*int` 

File size in bytes, if known. The true size is enforced
server-side during the upload. Maximum 5 GiB. Required for
`purpose: "avatar"` (maximum 10 MiB).
    
</dd>
</dl>

<dl>
<dd>

**purpose:** `*roamhq.CreateAssetRequestPurpose` 

What the asset will be used for. `file` (default) for chat
message attachments; `story` for story media (Personal tokens
only); `avatar` for `sender.imageUrl` and On-Air host photos.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Item
## Story
<details><summary><code>client.Story.Post(request) -> *roamhq.PostStoryResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Posts a story to your Roam. Stories are short photo or video updates that appear
above your profile picture for your teammates, and expire 24 hours after posting.

## Posting Flow

1. Create the media asset with [asset.create](https://developer.ro.am/docs/api/asset-create) using
   `purpose: "story"`, and upload the file bytes using the returned upload instructions.
2. Call this endpoint with the `assetId` (and an optional `caption`).

The media must be a photo or a video (videos up to 2.5 minutes; media is optimized
to portrait 1080×1920). If the upload is still processing — typical for videos in
the first seconds after upload — this endpoint returns a 400 with a "still
processing" message; retry after a short delay.

The media must outlive the story's 24-hour lifetime, so post within about 23 hours
of creating the asset (story assets expire about 48 hours after creation); older
assets are rejected and must be recreated.

**Access:** Personal only. Stories are always posted as the authenticated user —
a story appears above *your* profile picture, and there is no bot persona surface
for stories — so organization tokens are rejected.

**Required scope:** `chat:send_message` or `chat:write` (the same permission that
gates sending a chat message)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.PostStoryRequest{
    AssetID: "019be84b-0fa8-788f-8850-96de4cc39130",
    Caption: roamhq.String(
        "Greetings from the offsite 👋",
    ),
}
client.Story.Post(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**assetID:** `string` 

ID of a processed asset created via [asset.create](https://developer.ro.am/docs/api/asset-create)
with `purpose: "story"`. The asset must be owned by the authenticated user.
    
</dd>
</dl>

<dl>
<dd>

**caption:** `*string` — Optional caption displayed with the story.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## User
<details><summary><code>client.User.List() -> *roamhq.ListUserResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List workspace members, or hydrate an explicit ordered set of principal IDs.

Without `ids`, this is the active workspace member directory: guests,
bots, and archived/deactivated members are never enumerated. Members are
returned in the order they were added to the account.

With `ids`, the endpoint becomes an unpaginated principal hydrator. Pass one
comma-separated value containing at most 100 bare or tagged IDs. Duplicate
tokens are deduplicated in first-seen order; resolved entries are returned
in that order. Unknown IDs, groups, and unauthorized principals are silently
omitted. Explicit lookup may resolve archived/deactivated users and
authorized automated actors. The response keeps the existing `users` key
but its entries are principals, and `nextCursor` is omitted.

`ids` cannot be combined with `q`, `limit`, or `cursor`. `expand=status`
remains supported in either mode.

See [Identity & Principals](https://developer.ro.am/docs/guides/identity-and-principals).

**Required scope:** `user:read` (add `user:read.email` to include email addresses, `user:read.status` to expand presence status and `willReturn`)

**Access:** Organization and Personal.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListUserRequest{}
client.User.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ids:** `*string` 

One comma-separated list of up to 100 bare or tagged principal IDs.
Repeating the `ids` query parameter, including empty tokens, or combining
it with `q`, `limit`, or `cursor` returns `invalid_parameter`.
    
</dd>
</dl>

<dl>
<dd>

**q:** `*string` 

Case-insensitive member-directory filter by name. Also matches email
when the token has `user:read.email`. Cannot be combined with `ids`.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — The number of directory members to return per response. Default is 10. Cannot be combined with `ids`.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque directory cursor from a previous response's `nextCursor`. Cannot be combined with `ids`.
    
</dd>
</dl>

<dl>
<dd>

**expand:** `*string` — Comma-separated list of additional fields. Supported: `status` (requires `user:read.status`). Expanding `status` also returns `willReturn` when set.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.User.Info() -> *roamhq.User</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Resolve a v1 principal by ID, or look up a workspace member by email.

ID lookup resolves active or archived members, guests, and authorized
automated actors (classic bots, agents, assistants, and coworkers). The
response always includes `type: "user" | "bot"`; guests additionally have
`isGuest: true`. Groups, unknown IDs, and automated actors outside the
caller's Roam/account/owner boundary return `user_not_found`.

Email lookup remains workspace-member-only. Personal access tokens and the
MCP `user_info` tool may use ID lookup.

Provide either `id` or `email`, not both.

See [Identity & Principals](https://developer.ro.am/docs/guides/identity-and-principals) for the
taxonomy, visibility rules, and directory-versus-hydration guidance.

**Required scope:** `user:read` (add `user:read.email` to look up by email or include email in response, `user:read.status` to expand presence status, availability, and `willReturn`)

**Access:** Organization and Personal.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.InfoUserRequest{}
client.User.Info(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `*string` — A bare or tagged principal ID. Mutually exclusive with `email`.
    
</dd>
</dl>

<dl>
<dd>

**email:** `*string` — The user's email address. Mutually exclusive with `id`. Requires `user:read.email` scope.
    
</dd>
</dl>

<dl>
<dd>

**expand:** `*string` — Comma-separated list of additional fields to include. Supported: `status`, `available` (each requires `user:read.status`). Expanding `status` also returns `willReturn` when the user has a future out-of-office entry.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Users
<details><summary><code>client.Users.UserActivitySet(request) -> *roamhq.UserActivity</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Paint a badge (and optional glow) on a user's seat for work happening
outside Roam — a phone call, a browser meeting, a CRM session. Pass
`dnd: true` to also put their assigned office in Do Not Disturb.

The integration owns the lifecycle: `set` when the session starts,
`clear` when it ends. Re-posting the same `externalId` is the heartbeat
for long-running sessions — it refreshes `expiresAt` and, unless you
send `startedAt`, keeps the original start time. Roam stamps expiry
itself (default 10 minutes, maximum 60) so a dropped "ended" webhook
cannot leave a permanent glow.

`externalId` is unique per (integration, user). Two apps can hold
activities on the same person at once; you can only update or clear
your own rows.

See [External activity](https://developer.ro.am/docs/guides/user-activity) for display, DND,
TTL, stacking, and where the indicator appears on the map.

Identify the user with `userId`: a bare UUID, tagged `U-…` ID, or
ASCII email (same convention as `group.create` members). Third-party
systems that only have an email do not need a UUID lookup first.

**Access:** Organization and Personal. Organization tokens may target
any user in the workspace. Personal tokens (OAuth or PAT) may target
only the token owner.

**Required scope:** `user:write.activity`. Personal Access Tokens skip
this check; personal-mode OAuth installs must still request the scope.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.UserActivitySetRequest{
    UserID: "ada@example.com",
    ExternalID: "justcall:call:CA123",
    Display: &roamhq.UserActivityDisplay{
        Emoji: "📞",
        Title: "On a customer call",
        Subtitle: roamhq.String(
            "JustCall · Acme Corp",
        ),
        Color: roamhq.UserActivityDisplayColorGreen.Ptr(),
    },
    TTLSeconds: roamhq.Int(
        1800,
    ),
    Dnd: roamhq.Bool(
        true,
    ),
}
client.Users.UserActivitySet(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**userID:** `string` 

Target user. Bare UUID, tagged `U-…` ID, or ASCII email
(same convention as `group.create` members). Personal
tokens may only pass their own user. Does not require
`user:read.email` — email is an identifier, not a
disclosure.
    
</dd>
</dl>

<dl>
<dd>

**externalID:** `string` 

Caller-chosen session id, unique per integration and user.
Re-using it upserts the existing row (heartbeat). At most
128 Unicode code points.
    
</dd>
</dl>

<dl>
<dd>

**display:** `*roamhq.UserActivityDisplay` 
    
</dd>
</dl>

<dl>
<dd>

**ttlSeconds:** `*int` 

Seconds from now until expiry. Mutually exclusive with
`expiresAt`. Values above 3600 are **clamped** to 60
minutes, not rejected. Default when both are omitted: 600
(10 minutes).
    
</dd>
</dl>

<dl>
<dd>

**expiresAt:** `*time.Time` 

Absolute expiry (RFC3339, must be in the future). Mutually
exclusive with `ttlSeconds`. Instants more than 60 minutes
ahead are clamped to that maximum.
    
</dd>
</dl>

<dl>
<dd>

**startedAt:** `*time.Time` 

Optional session start (RFC3339). Omit on heartbeats to
preserve the original. A future value is clamped to the
server's now (clock skew; also so one integration cannot
pin the newest-first projection slot).
    
</dd>
</dl>

<dl>
<dd>

**dnd:** `*bool` 

If true, this activity contributes Do Not Disturb on the
user's **own assigned office** until it is cleared or
expires. Defaults to false — a badge does not lock an
office unless you opt in. Stacks with Zoom/Meet auto-DND
and other integrations' DND-flagged rows.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Users.UserActivityClear(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

End an activity previously created with [`user.activity.set`](https://developer.ro.am/docs/api/user-activity-set).
The row is keyed by this integration plus `userId` and `externalId` —
you cannot clear another app's activity.

Clearing a missing, already-cleared, or already-expired `externalId`
still returns **204**. Integrations retry "session ended" webhooks, and
the row may have expired in the meantime.

See [External activity](https://developer.ro.am/docs/guides/user-activity) for TTL, DND
stacking, and what happens on the map when the last activity clears.

**Access:** Organization and Personal. Organization tokens may target
any user in the workspace. Personal tokens (OAuth or PAT) may target
only the token owner.

**Required scope:** `user:write.activity`. Personal Access Tokens skip
this check; personal-mode OAuth installs must still request the scope.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.UserActivityClearRequest{
    UserID: "ada@example.com",
    ExternalID: "justcall:call:CA123",
}
client.Users.UserActivityClear(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**userID:** `string` 

Target user. Bare UUID, tagged `U-…` ID, or ASCII email
(same convention as `group.create` members). Personal
tokens may only pass their own user.
    
</dd>
</dl>

<dl>
<dd>

**externalID:** `string` — The `externalId` previously passed to `user.activity.set`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Users.UserActivityList() -> *roamhq.UserActivityListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Return every **currently live** external activity for a user — every
integration's rows, not only yours. Expired rows are omitted even
before the server reaper runs. Not paginated; ordered newest
`startedAt` first.

The map may show fewer entries than this list (the client projection
keeps the top three, always including at least one DND-flagged row).
`.list` is the source of truth for what is still live.

See [External activity](https://developer.ro.am/docs/guides/user-activity) for display, DND,
TTL, and where indicators appear.

**Access:** Organization and Personal. Organization tokens may list
any user in the workspace. Personal tokens (OAuth or PAT) may list
only the token owner.

**Required scope:** `user:read.activity`. Personal Access Tokens skip
this check; personal-mode OAuth installs must still request the scope.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.UserActivityListRequest{
    UserID: "userId",
}
client.Users.UserActivityList(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**userID:** `string` 

Target user. Bare UUID, tagged `U-…` ID, or ASCII email
(same convention as `group.create` members). Personal tokens
may only pass their own user.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Users.MessageeventExport(request) -> string</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Obtain a daily message event export containing DMs and group
chats within your account.

For customers with archival enabled (please reach out to a Roam
ArchiTech to get this process started), at the end of every day,
we export all message events for a particular day as a JSON Lines file.
This file contains all messages sent:
- by a Roam user who is a member of your organization
- into a chat containing (at the time of export) at least one Roam user who is a member of your organization
- by a bot integration that is part of your organization

This file also contains message edit and deletion events that meet the above criteria.
We specifically exclude waves, room invitations, and other non-message content
(that may appear as chats within the Roam application) from the export.

**Access:** Organization only.

**Required scope:** `admin:compliance:read`

### Message Event Structure

Each line within the file is a JSON object containing the following fields:
- eventType: a string that is one of “sent”, “edited”, or “deleted”
- chatId: a UUIDv4 identifier for a particular chat. All messages within the same chat shared the same chatId.
- threadTimestamp (optional): if part of a thread, the Unix epoch timestamp of the thread’s parent message in numerical format. All messages part of a thread share the same threadTimestamp.
- timestamp: the Unix epoch timestamp when the message was originally sent in numerical format.
- messageId: an internal UUIDv4 identifier as a string
- sender: a “Participant” object that identifiers the message sender
- contentType: a string that is one of the contentTypes associated with the “MessageContent” object
- content: a “MessageContent” object that contains the message’s content

### Participant

A Participant is a JSON object that contains three common fields: “participantType”, “id”, and “displayName”
- participantType: one of “email”, “bot”, or “occupant”
- id: a UUID identifier for the participant
- displayName: the name associated with the account or an empty string if not provided

Depending on the participant type, the object also contains additional fields:

Email Participant (a human user with a Roam user account)
- email: the email of the participant

Bot Participant (an automated user maintained by the Roam team or created via the Roam API)
- roamId: the roam ID associated with the integration
- integrationId: a unique integration ID name provided by the bot creator
- botCode: a unique identifier

### Message Content

A “MessageContent” object is a JSON object that contains the field “contentType” and,
depending on the content type, contains additional fields:

*Text Content* (contentType = “text”)
- text: the text in plaintext
- markdownText: the text in Markdown format
- attachments: A list of attachment objects

*Emoji Content* (contentType = “emoji”)
- text: text representation of the emoji
- colons: emoji in :emoji: format
- fileUrl: an optional field containing the URL to a custom emoji image

*Item Content* (contentType = “item”)
- itemUrl: the URL where the file can be downloaded from
- itemType: the type of item (e.g. "photo", "pdf", "blob", "video", "audio", etc.)

*Text Snippet Content* (contentType = "textSnippet")
- text: the content of the snippet
- language: the language of the snippet

*Members Changed Content* (contentType = “membersChanged”)
- added: a list of Participant objects corresponding to all participants added in this event
- removed: a list of Participant objects corresponding to all participants removed in this event
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.MessageeventExportRequest{
    Date: "2026-01-21",
}
client.Users.MessageeventExport(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**date:** `string` — The UTC date to fetch the export for in YYYY-MM-DD format.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## UserAuditLog
<details><summary><code>client.UserAuditLog.List() -> *roamhq.ListUserAuditLogResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get a list of user audit log entries for the account.

**Required scope:** `userauditlog:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListUserAuditLogRequest{}
client.UserAuditLog.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**date:** `*string` — The date to pull audit log entries from.  All activities from that date in UTC are returned.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Conversation
<details><summary><code>client.Conversation.List() -> *roamhq.ListConversationResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists conversations (meetings) that occurred in your Roam, with participant details.

**Access:**
- **Organization with [`admin:meetings:read`](https://developer.ro.am/docs/guides/scopes#meeting-width-adminmeetingsread)**
  (or a grandfathered roam-wide API key): all conversations in the workspace.
- **Personal access tokens:** supported — returns only conversations the
  token owner participated in (matched by confirmed email).
- **Organization without roam-wide meeting access** must use
  [`/meeting.list`](https://developer.ro.am/docs/api/meeting-list) instead (`403`).

**Required scope:** `meetings:read` (add `admin:meetings:read` for roam-wide org access)

Participant details require `user:read` scope. Email addresses require `user:read.email` scope.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListConversationRequest{}
client.Conversation.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**before:** `*time.Time` — Only return conversations that started before this ISO-8601 timestamp.
    
</dd>
</dl>

<dl>
<dd>

**after:** `*time.Time` — Only return conversations that started after this ISO-8601 timestamp.
    
</dd>
</dl>

<dl>
<dd>

**ascending:** `*bool` — Sort results in ascending order by start time. Default is descending (newest first).
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — The number of conversations to return per response.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Meeting
<details><summary><code>client.Meeting.List() -> *roamhq.ListMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List meetings, ordered newest-first.

**Access:** Organization and Personal. Personal tokens return meetings the
authenticated user participated in. Organization tokens return every meeting
in the Roam only with [`admin:meetings:read`](https://developer.ro.am/docs/guides/scopes#meeting-width-adminmeetingsread);
without it, results are limited to meetings the install's bot has access to.

**Required scope:** `meetings:read` (add `admin:meetings:read` for roam-wide org access)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListMeetingRequest{}
client.Meeting.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**before:** `*time.Time` — Only return meetings that started before this time (RFC-3339). Sub-millisecond precision is truncated.
    
</dd>
</dl>

<dl>
<dd>

**after:** `*time.Time` — Only return meetings that started after this time (RFC-3339). Sub-millisecond precision is truncated.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` 

Number of meetings to return per page. Capped to **10** when
`expand` includes `summary`, `actionItems`, or `chapters`, since
expanded payloads are substantially larger.
    
</dd>
</dl>

<dl>
<dd>

**expand:** `*string` 

Comma-separated list of fields to inline on each meeting. Allowed
values are `summary`, `actionItems`, and `chapters` — same shape
as on [`/meeting.info`](https://developer.ro.am/docs/api/meeting-info). Use this to
avoid N+1 follow-up calls when scanning many recent meetings.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.Info() -> *roamhq.InfoMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get detailed information about a specific meeting, including AI-generated summary, action items, and chapters.

Participants are included inline up to the `maxParticipants` limit. For meetings with more participants, use [`/meeting.participants`](https://developer.ro.am/docs/api/meeting-participants) to paginate through the full list.

**Access:** Organization and Personal. Personal tokens are limited to meetings
the authenticated user participated in. Organization tokens without
[`admin:meetings:read`](https://developer.ro.am/docs/guides/scopes#meeting-width-adminmeetingsread)
are limited to meetings the install's bot has access to.

**Required scope:** `meetings:read` (add `admin:meetings:read` for roam-wide org access; add `user:read` to include participants, `user:read.email` for participant emails)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.InfoMeetingRequest{
    ID: "id",
}
client.Meeting.Info(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The meeting ID.
    
</dd>
</dl>

<dl>
<dd>

**maxParticipants:** `*int` — Maximum number of participants to resolve and include inline. Use `/meeting.participants` for full pagination.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.Participants() -> *roamhq.ParticipantsMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Paginate through all participants of a meeting. This is the dedicated endpoint for retrieving the full participant list, complementing the capped inline participants in [`/meeting.info`](https://developer.ro.am/docs/api/meeting-info).

Pagination uses an **opaque cursor** (not a row offset). Pass `nextCursor`
from a previous response as `cursor` to fetch the next page. Invalid cursors
return `error: "invalid_cursor"` — see [Responses and Errors](https://developer.ro.am/docs/guides/responses-and-errors).

**Access:** Organization and Personal. Personal access tokens restrict to meetings the authenticated user participated in.

**Required scope:** `meetings:read` and `user:read` (add `user:read.email` for participant emails)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ParticipantsMeetingRequest{
    ID: "id",
}
client.Meeting.Participants(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The meeting ID.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Number of participants to return per page (default 50, max 200).
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not parse or construct cursors yourself.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.Transcript() -> *roamhq.TranscriptMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Retrieve the transcript for a meeting.

Supports content negotiation:
- **JSON** (default): Returns structured transcript with cues containing speaker IDs, text, and timing
- **WebVTT**: Set `Accept: text/vtt` header to receive standard WebVTT format with speaker voice tags

**Access:** Organization and Personal. Personal access tokens restrict to meetings the authenticated user participated in.

**Required scope:** `meetings:read`

**Errors** (see [Responses and Errors](https://developer.ro.am/docs/guides/responses-and-errors)):

| `error` code | Meaning |
|--------------|---------|
| `meeting_not_found` | Unknown or inaccessible meeting |
| `transcript_pending` | Not ready yet — retry later (may include `Retry-After`) |
| `transcript_unavailable` | Meeting was not transcribed — stop retrying |
| `upstream_timeout` | Timed out waiting on an upstream service — retry |
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.TranscriptMeetingRequest{
    ID: "id",
}
client.Meeting.Transcript(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The meeting ID.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.Search() -> *roamhq.SearchMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

AI-powered search across meeting transcripts and summaries.

**Access:** Personal access only. Organization (account-level) tokens are not supported.

**Required scope:** `meetings:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.SearchMeetingRequest{
    Query: "query",
}
client.Meeting.Search(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**query:** `string` — Search query string.
    
</dd>
</dl>

<dl>
<dd>

**after:** `*time.Time` — Only return results from meetings after this date (YYYY-MM-DD).
    
</dd>
</dl>

<dl>
<dd>

**before:** `*time.Time` — Only return results from meetings before this date (YYYY-MM-DD).
    
</dd>
</dl>

<dl>
<dd>

**timezone:** `*string` — Timezone for date interpretation (e.g. "America/New_York").
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.Prompt(request) -> *roamhq.PromptMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Ask an AI question about a meeting's transcript content. Returns a natural language response based on the meeting transcript.

**Access:** Organization and Personal. Personal access tokens restrict to meetings the authenticated user participated in.

**Required scope:** `meetings:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.PromptMeetingRequest{
    ID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    Prompt: "What action items were assigned to Alex?",
}
client.Meeting.Prompt(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The meeting ID.
    
</dd>
</dl>

<dl>
<dd>

**prompt:** `string` — The question to ask about the meeting.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.ShareLink(request) -> *roamhq.ShareLinkMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a shareable URL for a meeting that you can distribute to others. Pass the `id` of a meeting obtained from [`/meeting.list`](https://developer.ro.am/docs/api/meeting-list) or [`/meeting.info`](https://developer.ro.am/docs/api/meeting-info).

This endpoint is **get-or-create**: it returns the meeting's existing share link, or mints one the first time it is called for that meeting. Repeat calls for the same meeting return the same URL.

Creating a share link is a deliberate action, which is why it has its own endpoint rather than being returned as a field on `meeting.list` / `meeting.info` — fetching a meeting never mints a shareable link as a side effect. You can only create a share link for a meeting you can access; the same access check as `meeting.info` applies.

**Access:** Organization and Personal. Personal access tokens restrict to meetings the authenticated user participated in.

**Required scope:** `meetings:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ShareLinkMeetingRequest{
    ID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
}
client.Meeting.ShareLink(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The meeting ID.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.CreateLink(request) -> *roamhq.CreateLinkMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Create a meeting link.

**Access:** Organization and Personal. In Organization mode, specify the host by email. In Personal mode, the host defaults to the authenticated user.

**Required scope:** `meeting:write` or `meetinglink:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.CreateLinkMeetingRequest{
    Name: "Q1 Planning Session",
    Host: roamhq.String(
        "alex.chen@example.com",
    ),
    Start: roamhq.Time(
        roamhq.MustParseDateTime(
            "2026-02-15T14:00:00Z",
        ),
    ),
    End: roamhq.Time(
        roamhq.MustParseDateTime(
            "2026-02-15T15:00:00Z",
        ),
    ),
}
client.Meeting.CreateLink(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**name:** `string` — Meeting Name
    
</dd>
</dl>

<dl>
<dd>

**host:** `*string` 

Meeting Host Email, matching a member of your Roam.

Required for Organization tokens. For Personal tokens, this is optional and defaults to the authenticated user. If provided with a Personal token, it must match the authenticated user's email.
    
</dd>
</dl>

<dl>
<dd>

**start:** `*time.Time` — (Optional) Meeting start time in RFC3339.
    
</dd>
</dl>

<dl>
<dd>

**end:** `*time.Time` — (Optional) Meeting end time in RFC3339.
    
</dd>
</dl>

<dl>
<dd>

**requireUnconfirmedEmail:** `*bool` — (Optional) If true, guests must verify ownership of their email address before joining.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.LinkInfo(request) -> *roamhq.LinkInfoMeetingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get a meeting link.

**Access:** Organization and Personal. Personal tokens may only read meeting links where the authenticated user is the host.

**Required scope:** `meetinglink:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.LinkInfoMeetingRequest{
    ID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
}
client.Meeting.LinkInfo(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — Meeting Link ID
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Meeting.UpdateLink(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Update a meeting link.

**Access:** Organization and Personal. Personal tokens may only update meeting links where the authenticated user is the host.

**Required scope:** `meetinglink:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.UpdateLinkMeetingRequest{
    ID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    Name: "Q1 Planning Session - Updated",
    Start: roamhq.Time(
        roamhq.MustParseDateTime(
            "2026-02-15T15:00:00Z",
        ),
    ),
    End: roamhq.Time(
        roamhq.MustParseDateTime(
            "2026-02-15T16:30:00Z",
        ),
    ),
}
client.Meeting.UpdateLink(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — Meeting Link ID
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — Meeting Name
    
</dd>
</dl>

<dl>
<dd>

**host:** `*string` 

(Optional) Meeting Host Email.

The Host may NOT be updated.
As a result, this property may be omitted or empty.
If it is provided, it MUST match the existing value.
    
</dd>
</dl>

<dl>
<dd>

**start:** `*time.Time` — (Optional) Meeting start time in RFC3339.
    
</dd>
</dl>

<dl>
<dd>

**end:** `*time.Time` — (Optional) Meeting end time in RFC3339.
    
</dd>
</dl>

<dl>
<dd>

**requireUnconfirmedEmail:** `*bool` — (Optional) If true, guests must verify ownership of their email address before joining.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Meetings
<details><summary><code>client.Meetings.RecordingList() -> *roamhq.RecordingListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

**Legacy:** Prefer [`/meeting.list`](https://developer.ro.am/docs/api/meeting-list) /
[`/meeting.info`](https://developer.ro.am/docs/api/meeting-info) for new integrations.

Lists recordings in your home Roam, filtered by date range (after/before).
Organization clients without roam-wide meeting access
([`admin:meetings:read`](https://developer.ro.am/docs/guides/scopes#meeting-width-adminmeetingsread))
receive `403`; use [`/meeting.list`](https://developer.ro.am/docs/api/meeting-list) instead.
This route remains registered for existing callers. It returns v0-style
identifiers and is not a v1 media-download path.

The plural alias `/recordings.list` is also registered for existing callers;
use this singular form in new documentation and tooling.

The ordering of results depends on the filter specified:

- When no parameters are provided, the most recent recordings are returned,
  sorted in reverse chronological order. This is equivalent to specifying `before`
  as NOW and leaving `after` unspecified.

- If `after` is specified, the results are sorted in forward chronological order.

Either dates or datetimes may be specified. Dates are interpreted in UTC.

**Access:** Organization only. Requires roam-wide meeting access.

**Required scope:** `recordings:read` and `admin:meetings:read` (or a grandfathered roam-wide API key)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.RecordingListRequest{}
client.Meetings.RecordingList(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**after:** `*string` 

The datetime to begin listing recordings (YYYY-MM-DD or RFC-3339).
Defaults to "no filter".
    
</dd>
</dl>

<dl>
<dd>

**before:** `*string` 

The datetime until which to list recordings (YYYY-MM-DD or RFC-3339).
Defaults to "now".
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — The number of recordings to return per response. Default is 10.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Calendar
<details><summary><code>client.Calendar.CreateEvent(request) -> *roamhq.CreateEventCalendarResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Create a calendar event on the host's connected calendar. A Roam meeting link
is automatically attached and email notifications are sent to attendees.

The event is written to the first active, writable calendar associated with the
host. The host must have a connected calendar provider (e.g. Google, Microsoft).

**Recurring events:** Provide `rrule` to create a recurring series. A
`timeZone` is required for recurring events.

**All-day events:** Set `allDay: true`; `start` and `end` are interpreted as
dates and normalized to UTC midnight.

**Access:** Organization and Personal. For Organization tokens, the `host` email
is required and identifies the calendar owner. For Personal tokens, `host`
defaults to the authenticated user; if provided it must match the
authenticated user's email.

**Required scope:** `calendar:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.CreateEventCalendarRequest{
    Title: "Q1 Planning",
    Description: roamhq.String(
        "Plan Q1 roadmap",
    ),
    Start: roamhq.MustParseDateTime(
        "2026-02-15T14:00:00Z",
    ),
    End: roamhq.MustParseDateTime(
        "2026-02-15T15:00:00Z",
    ),
    TimeZone: roamhq.String(
        "America/Los_Angeles",
    ),
    Attendees: []string{
        "sam@example.com",
        "Alex Doe <alex@example.com>",
    },
    Host: roamhq.String(
        "host@example.com",
    ),
}
client.Calendar.CreateEvent(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**title:** `string` — Event title.
    
</dd>
</dl>

<dl>
<dd>

**description:** `*string` — (Optional) Event description.
    
</dd>
</dl>

<dl>
<dd>

**start:** `time.Time` — Event start time (RFC3339). For all-day events, the date portion is used.
    
</dd>
</dl>

<dl>
<dd>

**end:** `time.Time` — Event end time (RFC3339). For all-day events, the date portion is used.
    
</dd>
</dl>

<dl>
<dd>

**allDay:** `*bool` — Whether this is an all-day event. Defaults to false.
    
</dd>
</dl>

<dl>
<dd>

**rrule:** `*string` 

(Optional) iCalendar RFC 5545 recurrence rule, e.g. `FREQ=WEEKLY;COUNT=10`.
When provided, `timeZone` is required.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `*string` 

IANA timezone name, e.g. `America/New_York`. Required for recurring
events; recommended for all events. Defaults to `UTC` when omitted.
    
</dd>
</dl>

<dl>
<dd>

**attendees:** `[]string` 

Attendee email addresses. Each entry may be a plain email
(`user@example.com`) or an address string (`Name <user@example.com>`).
    
</dd>
</dl>

<dl>
<dd>

**host:** `*string` 

Calendar host email. Required for Organization tokens. For Personal
tokens, defaults to the authenticated user and, if provided, must
match the authenticated user's email.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Calendar.List() -> *roamhq.ListCalendarResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List events from the authenticated user's connected calendars within
a date range.

Pulls events from every active personal calendar attached to the user
(e.g. Google, Microsoft) and merges them into a single chronological
list. Canceled events are omitted.

**Date range:** Defaults to a 7-day window starting today (caller's
timezone). Pass `startDate` to shift the window's start; pass
`endDate` to set its end (inclusive). Both are interpreted as
`YYYY-MM-DD` in the caller's timezone.

**Access:** Personal access only. Organization tokens do not have
access to individual calendars and receive a `400`.

**Required scope:** `calendar:read`

`meetings:read` also grants this endpoint, but only for API clients
registered **before 2026-07-29T00:00Z**. Clients registered on or after that
date must hold `calendar:read`, or the call fails with `403` /
`missing_scope`. See [Scopes](https://developer.ro.am/docs/guides/scopes).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListCalendarRequest{}
client.Calendar.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**startDate:** `*string` — First day to include (`YYYY-MM-DD`, caller's timezone). Defaults to today.
    
</dd>
</dl>

<dl>
<dd>

**endDate:** `*string` 

Last day to include (`YYYY-MM-DD`, caller's timezone, inclusive).
Defaults to seven days after the resolved `startDate`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Lobby
<details><summary><code>client.Lobby.List() -> *roamhq.ListLobbyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists active lobbies in your account.

A lobby URL has the form `ro.am/{handle}` or `ro.am/{handle}/{slug}`.
- The "handle" is the first path segment
- The "slug" is the optional second path segment. It may be empty for the default lobby under a handle

Optionally filter by a specific lobby handle. If provided, only lobbies
associated with that handle are returned.

This endpoint is **not paginated**. The 200 body is `{ "lobbies": [...] }`
with every matching lobby; there is no `cursor` / `nextCursor` and no
`data` array. The TypeScript SDK returns that object directly, not a
page helper.

**Access:** Organization and Personal.

**Required scope:** `lobby:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListLobbyRequest{}
client.Lobby.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**handle:** `*string` 

Filter by lobby handle (first path segment), e.g., `robfig` for
`ro.am/robfig` or `ro.am/robfig/tour`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Lobby.ListBookings() -> *roamhq.ListBookingsLobbyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists bookings for a specific lobby configuration, filtered by date range (after/before).

The ordering of results depends on the filter specified:

- When no parameters are provided, the most recent bookings are returned,
  sorted in reverse chronological order. This is equivalent to specifying `before`
  as NOW and leaving `after` unspecified.

- If `after` is specified, the results are sorted in forward chronological order.

Either dates or datetimes may be specified. Dates are interpreted in UTC.

**Access:** Organization and Personal.

**Required scope:** `lobby:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListBookingsLobbyRequest{
    LobbyID: "lobbyId",
}
client.Lobby.ListBookings(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**lobbyID:** `string` — The lobby configuration ID to list bookings for.
    
</dd>
</dl>

<dl>
<dd>

**after:** `*time.Time` 

The datetime to begin listing bookings (YYYY-MM-DD or RFC-3339).
Defaults to "no filter".
    
</dd>
</dl>

<dl>
<dd>

**before:** `*time.Time` 

The datetime until which to list bookings (YYYY-MM-DD or RFC-3339).
Defaults to "now".
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — The number of bookings to return per response. Default is 10.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Magicast
<details><summary><code>client.Magicast.List() -> *roamhq.ListMagicastResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List Magicasts in your account, most recent first.

Returns metadata only (`id`, `name`, `createdAt`, `ownerId`,
`coverImageUrl`). Use [`/magicast.info`](https://developer.ro.am/docs/api/magicast-info) for
transcript cues, chapters, video status, and a signed download URL.

**Access:** Organization and Personal. Organization tokens list every
Magicast in the account, including ones the creator never shared. Personal
tokens are restricted to Magicasts owned by the authenticated user.

**Required scope:** `magicast:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListMagicastRequest{}
client.Magicast.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**after:** `*time.Time` — Only return magicasts created after this time (RFC-3339).
    
</dd>
</dl>

<dl>
<dd>

**before:** `*time.Time` — Only return magicasts created before this time (RFC-3339).
    
</dd>
</dl>

<dl>
<dd>

**ascending:** `*bool` — Sort oldest-first instead of newest-first.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Number of magicasts to return per response. Default 10.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Magicast.Info() -> *roamhq.MagicastInfo</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get details for a single Magicast by ID, including transcript cues,
chapters, video status, a signed video download URL when ready, and a
player URL if a share link already exists.

This is the content endpoint. [`/magicast.list`](https://developer.ro.am/docs/api/magicast-list)
returns metadata only. Magicasts are not meetings — they do not appear on
[`/recording.list`](https://developer.ro.am/docs/api/recording-list) or meeting transcript
surfaces, and they have no Magic Minutes summary or action items.

Asset, transcript, and share-link lookups are best-effort. If the video or
transcript is still processing, those fields are omitted and the request
still succeeds. Fetching this endpoint **never** mints a shareable link;
use [`/magicast.shareLink`](https://developer.ro.am/docs/api/magicast-share-link) for that.

There is no `https://ro.am/magicast/{id}` browser URL. The player URL is
always `https://ro.am/share/{key}`.

**Access:** Organization and Personal. Organization tokens can read every
Magicast in the account, including ones the creator never shared. Personal
tokens are restricted to Magicasts owned by the authenticated user. Filter
on whether `shareUrl` is present if you only want shared recordings.

**Required scope:** `magicast:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.InfoMagicastRequest{
    ID: "id",
}
client.Magicast.Info(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The magicast ID.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Magicasts
<details><summary><code>client.Magicasts.MagicastShareLink(request) -> *roamhq.MagicastShareLinkResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a shareable player URL for a Magicast. Pass the `id` obtained from
[`/magicast.list`](https://developer.ro.am/docs/api/magicast-list) or
[`/magicast.info`](https://developer.ro.am/docs/api/magicast-info).

This endpoint is **get-or-create**: it returns the Magicast's existing
share link, or mints one the first time it is called. Repeat calls for the
same Magicast return the same URL.

Creating a share link is a deliberate action, which is why it has its own
endpoint rather than being returned as a field that is always present on
`magicast.list` / `magicast.info`. Fetching a Magicast never mints a
shareable link as a side effect. `magicast.info` includes `shareUrl` only
when a link already exists.

The URL is `https://ro.am/share/{key}`. There is no
`https://ro.am/magicast/{id}` route.

You can only create a share link for a Magicast you can access; the same
access check as `magicast.info` applies.

**Access:** Organization and Personal. Personal access tokens restrict to
Magicasts owned by the authenticated user.

**Required scope:** `magicast:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.MagicastShareLinkRequest{
    ID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
}
client.Magicasts.MagicastShareLink(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The Magicast ID.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Group
<details><summary><code>client.Group.List() -> *roamhq.ListGroupResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists non-archived groups accessible to the caller.

Filter by name with `query` (ranked text match), restrict by group
type with `type`, and paginate with `limit` / `cursor`.

**Access:** Organization and Personal.

**Required scope:** `group:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ListGroupRequest{}
client.Group.List(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**query:** `*string` — Text filter. Groups are ranked by how well their name matches the query.
    
</dd>
</dl>

<dl>
<dd>

**type_:** `*string` 

Comma-separated list of group types to include. Must be one or
more of `standard`, `magicast`, `meeting`, `roam`, `onair`.
Defaults to all types.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Number of groups to return per page (default 50, max 100).
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Info() -> *roamhq.Group</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get information about a specific group by its ID or name.

Provide either `id` or `name`, not both.

**Required scope:** `group:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.InfoGroupRequest{}
client.Group.Info(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `*string` — The group's ID. Mutually exclusive with `name`.
    
</dd>
</dl>

<dl>
<dd>

**name:** `*string` — The group's name. Mutually exclusive with `id`. Returns first match if multiple groups have the same name.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Create(request) -> *roamhq.Group</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Create a group chat.

Groups which specify at least one admin will operate in an "Admin only" management
mode, where only admins may change settings. Otherwise, all members have
that capability.

Groups require at least one member. Users can be specified by user ID or email address.
Unrecognized emails are invited as group members only — they do not receive a
[Guest Badge](https://developer.ro.am/docs/guides/guest-badges) unless you also call
[`guest.badge.create`](https://developer.ro.am/docs/api/guest-badge-create).

**Required scope:** `group:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.CreateGroupRequest{
    Name: "Engineering Team",
    Description: roamhq.String(
        "Group chat for engineering discussions and updates",
    ),
    Private: roamhq.Bool(
        false,
    ),
    EnforceThreads: roamhq.Bool(
        true,
    ),
    Members: []*roamhq.CreateGroupRequestMembersItem{
        &roamhq.CreateGroupRequestMembersItem{
            UserID: "alex.chen@example.com",
            Role: roamhq.CreateGroupRequestMembersItemRoleMember,
        },
        &roamhq.CreateGroupRequestMembersItem{
            UserID: "taylor@example.com",
            Role: roamhq.CreateGroupRequestMembersItemRoleMember,
        },
        &roamhq.CreateGroupRequestMembersItem{
            UserID: "jordan.smith@example.com",
            Role: roamhq.CreateGroupRequestMembersItemRoleAdmin,
        },
    },
}
client.Group.Create(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**name:** `string` — Name of the group
    
</dd>
</dl>

<dl>
<dd>

**description:** `*string` — Description of the group
    
</dd>
</dl>

<dl>
<dd>

**private:** `*bool` — Whether the group is private (default false)
    
</dd>
</dl>

<dl>
<dd>

**enforceThreads:** `*bool` — Whether to enforce threaded conversations
    
</dd>
</dl>

<dl>
<dd>

**members:** `[]*roamhq.CreateGroupRequestMembersItem` — Group members with their roles
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Rename(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Rename a group by ID.

Apps may only rename groups for which they are an admin.

**Required scope:** `group:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.RenameGroupRequest{
    ID: "88bebce7-6cbb-4666-96f9-5c02d73e6661",
    Name: "Product Engineering",
}
client.Group.Rename(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The group ID
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — The new name for the group
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Archive(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Archive a group by ID.

Apps may only archive groups for which they are an admin.

**Required scope:** `group:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.ArchiveGroupRequest{
    ID: "88bebce7-6cbb-4666-96f9-5c02d73e6661",
}
client.Group.Archive(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The group ID to archive.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Members() -> *roamhq.MembersGroupResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List members in a group with their roles.

Apps may list members if one of the following conditions is true:
1. It is a public group in their Roam.
2. They are a member of the group.

**Required scope:** `group:read`

Every returned `userId` is a visible principal ID that resolves through
[`user.info`](https://developer.ro.am/docs/api/user-info) with the same credentials. Use
`user.list?ids` for ordered bulk hydration.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.MembersGroupRequest{
    ID: "id",
}
client.Group.Members(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — Group ID.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — The number of members to return per response. Default is 10.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Add(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Add one or more group members with specified roles.

Members can be specified by user ID or email address. Each member must be assigned a role (member or admin).

Adding an unrecognized email does **not** grant a [Guest Badge](https://developer.ro.am/docs/guides/guest-badges). Use [`guest.badge.create`](https://developer.ro.am/docs/api/guest-badge-create) first if the person is not a workspace member.

Apps may add members to a group if one of the following conditions is true:
1. It is a public group in their Roam.
2. They are a member of the group.

If attempting to add an admin, the app must be an admin of the group.

**Required scope:** `group:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.AddGroupRequest{
    ID: "88bebce7-6cbb-4666-96f9-5c02d73e6661",
    Members: []*roamhq.AddGroupRequestMembersItem{
        &roamhq.AddGroupRequestMembersItem{
            UserID: "709b8a57-70bc-427a-b6f0-b16ba5297f8c",
            Role: roamhq.AddGroupRequestMembersItemRoleMember,
        },
        &roamhq.AddGroupRequestMembersItem{
            UserID: "f589a8cb-78ac-493e-8719-0fa8a22f65e0",
            Role: roamhq.AddGroupRequestMembersItemRoleMember,
        },
        &roamhq.AddGroupRequestMembersItem{
            UserID: "af6663d5-0f37-4105-95df-4fea20ef7c7c",
            Role: roamhq.AddGroupRequestMembersItemRoleAdmin,
        },
    },
}
client.Group.Add(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — Group ID
    
</dd>
</dl>

<dl>
<dd>

**members:** `[]*roamhq.AddGroupRequestMembersItem` — List of members to add with their roles
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Join(request) -> *roamhq.Group</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Join a public group as the calling identity (Slack `conversations.join`).

- Org tokens add the bot address as a member.
- Personal tokens add the **owner person**, never the PAT bot address.
- Private groups cannot be self-joined (`403`).
- Idempotent if the calling identity is already a member.
- Non-members of a group in another roam receive an opaque `403`
  (`group_not_found`) — archived / type / privacy are not distinguished.

Why join (webhooks vs history vs post): [Chat](https://developer.ro.am/docs/guides/chat).

**Access:** Organization and Personal.

**Required scope:** `group:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.JoinGroupRequest{
    ID: "88bebce7-6cbb-4666-96f9-5c02d73e6661",
}
client.Group.Join(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — Group ID
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Group.Remove(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove one or more group members.

Members can be specified by user ID or email address.

Apps may remove members from a group if one of the following conditions is true:
1. It is a public group in their Roam.
2. They are a member of the group.

Removing members with the Admin role is not yet supported.

**Required scope:** `group:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.RemoveGroupRequest{
    ID: "88bebce7-6cbb-4666-96f9-5c02d73e6661",
    Members: []string{
        "709b8a57-70bc-427a-b6f0-b16ba5297f8c",
    },
}
client.Group.Remove(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — Group ID
    
</dd>
</dl>

<dl>
<dd>

**members:** `[]string` — List of member IDs or email addresses to remove
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Groups
<details><summary><code>client.Groups.List() -> []*roamhq.GroupsListResponseItem</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

**Legacy:** Prefer [`/group.list`](https://developer.ro.am/docs/api/group-list) for new integrations.

Lists all public, non-archived groups in your home Roam.

Unlike `/group.list`, this endpoint returns a **raw JSON array** (not the
`{"ok": true, …}` envelope). It is the sole ok-envelope exception on `/v1`
and remains only for existing callers.

**Access:** Organization only.

**Required scope:** `group:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Groups.List(
    context.TODO(),
)
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Guest Badges
<details><summary><code>client.GuestBadges.GuestBadgeCreate(request) -> *roamhq.GuestBadge</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Grant a Guest Badge so an email that is **not** a workspace member can
visit a host in Roam.

This is **not** an [On-Air event guest](https://developer.ro.am/docs/onair-api/on-air-api). It is
also **not** implied by [`group.add`](https://developer.ro.am/docs/api/group-add): adding an email
to a group does not mint a badge or send the invite. Typical onboarding is
`guest.badge.create` then `group.add`.

Repeating create for the same host and email returns the existing badge
(`visitPermission` is **not** updated) and does not re-send the invite.
To flip on-map access after create, use
[`guest.badge.update`](https://developer.ro.am/docs/api/guest-badge-update).

**Access:** Organization and Personal.
Organization tokens require `hostUserId`. Personal tokens default to the
token owner; naming a different host returns `403` `access_mode_not_supported`.

**Required scope:** `guest:write`. Personal Access Tokens use the
`pat:guests:write` group.

See [Guest Badges](https://developer.ro.am/docs/guides/guest-badges).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.GuestBadgeCreateRequest{
    Email: "alex@client.example",
    HostUserID: roamhq.String(
        "3f1c0b2a-8d4e-4c91-9a7b-2e6f1d8c0a11",
    ),
}
client.GuestBadges.GuestBadgeCreate(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**email:** `string` — Guest email. ASCII only. Must not be a workspace member.
    
</dd>
</dl>

<dl>
<dd>

**hostUserID:** `*string` 

Host member. UUID or member email — the same convention as
[`group.create`](https://developer.ro.am/docs/api/group-create) `members[].userId`.
Required for organization tokens. Optional for personal tokens
(defaults to the token owner).
    
</dd>
</dl>

<dl>
<dd>

**visitPermission:** `*bool` 

Whether the guest may visit the host on the map. Defaults to
`true`. Ignored on an idempotent retry of an existing grant
— use [`guest.badge.update`](https://developer.ro.am/docs/api/guest-badge-update)
to change it.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.GuestBadges.GuestBadgeList() -> *roamhq.GuestBadgeListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List issued Guest Badges.

Organization tokens return every issued badge in the workspace. Personal
tokens return only badges the token owner issued. This is the issued
(host) view — the same rows `guest.badge.create` returns — not the
guest's hidden-inbox view.

Filter with `email` (alias-aware) and/or `hostUserId` (UUID or member
email). Paginate with `limit` / `cursor` (default 50, max 100). Results
are sorted by `(hostUserId, email)`.

Organization keys that only have `guest:write` must also request
`guest:read` to call list. Personal Access Tokens with `pat:guests:write`
already include `guest:read`.

**Access:** Organization and Personal.
Personal tokens naming a different host return `403`
`access_mode_not_supported`.

**Required scope:** `guest:read`. Personal Access Tokens use the
`pat:guests:write` group.

See [Guest Badges](https://developer.ro.am/docs/guides/guest-badges).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.GuestBadgeListRequest{}
client.GuestBadges.GuestBadgeList(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**email:** `*string` 

Guest email. ASCII only. Matches the stored address and its verified
domain aliases.
    
</dd>
</dl>

<dl>
<dd>

**hostUserID:** `*string` 

Host member. UUID or member email — the same convention as
[`group.create`](https://developer.ro.am/docs/api/group-create) `members[].userId`.
Archived hosts may be named (they typically have no remaining grants).
Personal tokens may only pass the token owner.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Number of badges to return per page (default 50, max 100).
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`. Do not construct cursors manually.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.GuestBadges.GuestBadgeUpdate(request) -> *roamhq.GuestBadge</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Update `visitPermission` on an existing Guest Badge.

[`guest.badge.create`](https://developer.ro.am/docs/api/guest-badge-create) is idempotent and
does **not** change `visitPermission` on an existing grant. Use this
endpoint to flip on-map visit access after create.

If only one host in the workspace has granted this email, `hostUserId`
may be omitted. If several hosts have, pass `hostUserId` to pick which
grant to update (`400` `missing_parameter` otherwise). Same-host alias
rows are updated together.

Personal tokens can only update badges they issued. Naming another host
is `403` `access_mode_not_supported`; omitting `hostUserId` when the
token owner has no matching grant is `404` `not_found`.

**Access:** Organization and Personal.

**Required scope:** `guest:write`. Personal Access Tokens use the
`pat:guests:write` group.

See [Guest Badges](https://developer.ro.am/docs/guides/guest-badges).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.GuestBadgeUpdateRequest{
    Email: "alex@client.example",
    HostUserID: roamhq.String(
        "3f1c0b2a-8d4e-4c91-9a7b-2e6f1d8c0a11",
    ),
    VisitPermission: false,
}
client.GuestBadges.GuestBadgeUpdate(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**email:** `string` — Guest email. ASCII only.
    
</dd>
</dl>

<dl>
<dd>

**hostUserID:** `*string` 

Host member. UUID or member email. Required when more than one
host has granted this email. Optional for a unique grant, and
for personal tokens (defaults to the token owner).
    
</dd>
</dl>

<dl>
<dd>

**visitPermission:** `bool` — Whether the guest may visit the host on the map.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.GuestBadges.GuestBadgeRevoke(request) -> *roamhq.GuestBadgeRevokeResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Revoke Guest Badge(s) for an email.

If only one host in the workspace has granted this email, `hostUserId` may
be omitted. If several hosts have, pass `hostUserId` to pick which grant
to revoke (`400` `missing_parameter` otherwise). Same-host alias rows are
all revoked together.

Returns `{ "revoked": true }` when a matching grant was found and deleted,
or `{ "revoked": false }` when there was nothing to revoke (already gone,
including after the host was archived — archiving a member deletes the
badges they granted). Naming an archived host does not 404.

Personal tokens can only revoke badges they issued. Naming another host is
`403` `access_mode_not_supported`; omitting `hostUserId` when only another
host granted the email is a no-op (`revoked: false`).

**Access:** Organization and Personal.

**Required scope:** `guest:write`. Personal Access Tokens use the
`pat:guests:write` group.

See [Guest Badges](https://developer.ro.am/docs/guides/guest-badges).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.GuestBadgeRevokeRequest{
    Email: "alex@client.example",
}
client.GuestBadges.GuestBadgeRevoke(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**email:** `string` — Guest email. ASCII only.
    
</dd>
</dl>

<dl>
<dd>

**hostUserID:** `*string` 

Host member. UUID or member email. Required when more than one
host has granted this email. Optional for a unique grant, and
for personal tokens (defaults to the token owner).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Token
<details><summary><code>client.Token.Info() -> *roamhq.InfoTokenResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get information about the access token, including the authenticated user/bot
and granted scopes.

**No specific scope required.**
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Token.Info(
    context.TODO(),
)
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Token.Revoke() -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Permanently revoke the presented OAuth access token **and its refresh
token**. After a successful response the grant is dead — refresh will not
resurrect it; the client must re-authorize.

This does **not** uninstall your app from the workspace, delete webhook
subscriptions, or affect other users' tokens. For install removal see the
[`app.uninstalled`](https://developer.ro.am/docs/webhooks/app-uninstalled) event (fired from admin
/ Dev Settings uninstall paths, not from this endpoint).

On success Roam also delivers a [`token.revoked`](https://developer.ro.am/docs/webhooks/token-revoked)
webhook to your subscriptions (`reason: "api_revoked"`), including to the
same app that called this endpoint. Treat that delivery as idempotent.

Subsequent API calls with the revoked access token return HTTP `401` with
`invalid_token` (the token row is gone). Distinct from `token_revoked`,
which signals an archived person or archived client while a credential may
still exist.

This operation is only valid for OAuth access tokens, not for API keys.

**Access:** Organization and Personal (OAuth access tokens only). Personal
tokens may revoke their own grant. API keys cannot use this endpoint.

**No specific scope required.**
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Token.Revoke(
    context.TODO(),
)
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Webhook
<details><summary><code>client.Webhook.List() -> *roamhq.ListWebhookResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List all webhook subscriptions owned by the authenticated API client.

The response includes both **dynamic** subscriptions (created via
[`/webhook.subscribe`](https://developer.ro.am/docs/webhooks/webhook-subscribe)) and **static**
subscriptions configured in the Roam Administration UI.

Each object may include `lastSuccessAt`, `failStreakStartedAt`, and
`disabledAt` (omitted when null). `disabledAt` means the destination is
paused. See [Subscription health](https://developer.ro.am/docs/webhooks/webhooks#subscription-health).

**Required scope:** `webhook:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Webhook.List(
    context.TODO(),
)
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhook.Subscribe(request) -> *roamhq.Webhook</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Create or update a webhook subscription for a given event. If a subscription
already exists for the same event and URL, its filter is updated instead of
creating a duplicate. Re-subscribing the same event and URL also clears a
pause (`disabledAt` / `failStreakStartedAt`) so deliveries resume on the
next event. See [Subscription health](https://developer.ro.am/docs/webhooks/webhooks#subscription-health).

**Event names are dotted:** `chat.message`, `lobby.booked`,
`magicast.created`. Colon names (`chat:message:dm`, `lobby:booked`) are
v0-only — sending them here returns `400` / `Unrecognized event`.

Roam does not probe the destination URL when you subscribe — the
subscription is created immediately and the first delivery is a real event.

Optional `filter` limits which occurrences are delivered. Which keys are
valid depends on `event` — see that event's page and the
[Event Filters](https://developer.ro.am/docs/webhooks/webhooks#event-filters) table. Omit
`filter` to receive every occurrence. An empty object (`{}`) is rejected,
as is a filter that does not apply to the event.

**DMs only:**

```json
{
  "url": "https://example.com/hooks/messages",
  "event": "chat.message",
  "filter": { "chatType": "dm" }
}
```

**Required scope:** `webhook:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.WebhookSubscriptionRequest{
    URL: "https://example.com/hooks/messages",
    Event: roamhq.WebhookSubscriptionRequestEventChatMessage,
}
client.Webhook.Subscribe(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**url:** `string` — Destination URL for webhook deliveries (max 1024 characters). HTTPS is required outside local environments.
    
</dd>
</dl>

<dl>
<dd>

**event:** `*roamhq.WebhookSubscriptionRequestEvent` — Event to subscribe to.
    
</dd>
</dl>

<dl>
<dd>

**filter:** `*roamhq.WebhookSubscriptionFilter` 

Optional event-specific filter. Which keys are valid depends on `event`
(see the schema). Omit to receive every occurrence; `{}` and `null` are
rejected rather than treated as "omitted". Example for DMs only:
`{"chatType": "dm"}`.
    
</dd>
</dl>

<dl>
<dd>

**apiVersion:** `*string` 

Optional [API version](https://developer.ro.am/docs/guides/api-versioning) (`YYYY-MM-DD`) to pin
this subscription's payload shape to. When omitted, the subscription is
frozen at your integration's default version. Unsupported values return
`400`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhook.Unsubscribe(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove a webhook subscription by ID.

The request body is JSON: `{"id": "<subscription uuid>"}`. This differs
from v0, which expects `application/x-www-form-urlencoded` with the same
`id` field. Sending JSON to `/v0/webhook.unsubscribe` returns
`id parameter required`.

**Required scope:** `webhook:write`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.WebhookUnsubscribeRequest{
    ID: "19c6401f-6d02-4d8c-87c5-9fc45f02f4b5",
}
client.Webhook.Unsubscribe(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — Identifier of the webhook subscription to remove.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhook.Deliveries() -> *roamhq.DeliveriesWebhookResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List recent **failed** webhook delivery attempts for the authenticated API
client, newest first. Use this to debug a misbehaving endpoint and to find
the events you need to replay: successful (2xx) deliveries are never
recorded, so every row here is a delivery your endpoint did not accept.

Timeouts are first-class failures: `statusCode` is `0` and `error` is
`timeout`. For HTTP error responses, a truncated copy of your server's
response body is included to aid debugging. The request payload is never
stored — to recover the data, re-fetch the underlying resource (e.g. via
`chat.history`) using the delivery's `messageId`/`event` context.

Results are strictly scoped to the caller's own subscriptions and retained
for roughly 30 days.

**Access:** Organization and Personal.

**Required scope:** `webhook:read`
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &roamhq.DeliveriesWebhookRequest{}
client.Webhook.Deliveries(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**webhook:** `*string` — Only return deliveries for this webhook subscription ID.
    
</dd>
</dl>

<dl>
<dd>

**event:** `*string` — Only return deliveries for this event name (e.g. `chat.message`).
    
</dd>
</dl>

<dl>
<dd>

**after:** `*string` — Only return deliveries after this time (RFC3339 or `YYYY-MM-DD`). Results switch to oldest-first.
    
</dd>
</dl>

<dl>
<dd>

**before:** `*string` — Only return deliveries before this time (RFC3339 or `YYYY-MM-DD`).
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Maximum number of deliveries to return.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque pagination cursor from a previous response's `nextCursor`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

