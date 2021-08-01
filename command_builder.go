package router

import "github.com/Postcord/objects"

type commandBuilder struct {
	name string
	map_ map[string]interface{}
	cmd  Command
}

func (c *commandBuilder) Description(description string) CommandBuilder {
	c.cmd.Description = description
	return c
}

func (c *commandBuilder) Option(option *objects.ApplicationCommandOption) CommandBuilder {
	c.cmd.Options = append(c.cmd.Options, option)
	return c
}

func (c *commandBuilder) DefaultPermission() CommandBuilder {
	c.cmd.DefaultPermission = true
	return c
}

func (c *commandBuilder) AllowedMentions(config *objects.AllowedMentions) CommandBuilder {
	c.cmd.AllowedMentions = config
	return c
}

func (c *commandBuilder) Handler(handler func(*CommandRouterCtx) error) CommandBuilder {
	c.cmd.Function = handler
	return c
}

func (c *commandBuilder) Build() (*Command, error) {
	c.map_[c.name] = &c.cmd
	return &c.cmd, nil
}

func (c *commandBuilder) MustBuild() *Command {
	cmd, err := c.Build()
	if err != nil {
		panic(err)
	}
	return cmd
}

func (c textCommandBuilder) Description(description string) TextCommandBuilder {
	c.commandBuilder.Description(description)
	return c
}

func (c textCommandBuilder) Handler(handler func(*CommandRouterCtx) error) TextCommandBuilder {
	c.commandBuilder.Handler(handler)
	return c
}

func (c messageCommandBuilder) Handler(handler func(*CommandRouterCtx, *objects.Message) error) MessageCommandBuilder {
	c.commandBuilder.Handler(messageTargetWrapper(handler))
	return c
}

func (c userCommandBuilder) Handler(handler func(*CommandRouterCtx, *objects.GuildMember) error) UserCommandBuilder {
	c.commandBuilder.Handler(memberTargetWrapper(handler))
	return c
}

// TextCommandBuilder is used to define a builder for a Command object where the type is a text command.
type TextCommandBuilder interface {
	// Description is used to define the commands description.
	Description(string) TextCommandBuilder

	// Option is used to add a command option.
	Option(*objects.ApplicationCommandOption) TextCommandBuilder

	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() TextCommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) TextCommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx) error) TextCommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// MessageCommandBuilder is used to define a builder for a Message object where the type is a user command.
type MessageCommandBuilder interface {
	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() MessageCommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) MessageCommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx, *objects.Message) error) MessageCommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// UserCommandBuilder is used to define a builder for a Command object where the type is a user command.
type UserCommandBuilder interface {
	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() UserCommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) UserCommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx, *objects.GuildMember) error) UserCommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// CommandBuilder is used to define a builder for a Command object where the type isn't known.
type CommandBuilder interface {
	// Description is used to define the commands description.
	Description(string) CommandBuilder

	// TextCommand is used to define that this should be a text command builder.
	TextCommand() TextCommandBuilder

	// MessageCommand is used to define that this should be a message command builder.
	MessageCommand() MessageCommandBuilder

	// UserCommand is used to define that this should be a message command builder.
	UserCommand() UserCommandBuilder

	// Option is used to add a command option.
	Option(*objects.ApplicationCommandOption) CommandBuilder

	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() CommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) CommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx) error) CommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// NewCommandBuilder is used to create a builder for a *Command object.
func (c CommandGroup) NewCommandBuilder(name string) TextCommandBuilder {
	x := &commandBuilder{name: name, map_: c.Subcommands, cmd: Command{commandType: int(objects.CommandTypeChatInput)}}
	return x.TextCommand()
}
