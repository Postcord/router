// Code generated by generate_command_builder.go; DO NOT EDIT.

package router

//go:generate go run generate_command_builder.go

import "github.com/Postcord/objects"

type textCommandBuilder struct {
	*commandBuilder
}

func (c textCommandBuilder) Option(option *objects.ApplicationCommandOption) TextCommandBuilder {
	c.commandBuilder.Option(option)
	return c
}

func (c textCommandBuilder) DefaultPermission() TextCommandBuilder {
	c.commandBuilder.DefaultPermission()
	return c
}

func (c textCommandBuilder) AllowedMentions(config *objects.AllowedMentions) TextCommandBuilder {
	c.commandBuilder.AllowedMentions(config)
	return c
}

func (c *commandBuilder) TextCommand() TextCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeChatInput)
	return textCommandBuilder{c}
}

type messageCommandBuilder struct {
	*commandBuilder
}

func (c messageCommandBuilder) Option(option *objects.ApplicationCommandOption) MessageCommandBuilder {
	c.commandBuilder.Option(option)
	return c
}

func (c messageCommandBuilder) DefaultPermission() MessageCommandBuilder {
	c.commandBuilder.DefaultPermission()
	return c
}

func (c messageCommandBuilder) AllowedMentions(config *objects.AllowedMentions) MessageCommandBuilder {
	c.commandBuilder.AllowedMentions(config)
	return c
}

func (c *commandBuilder) MessageCommand() MessageCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeMessage)
	return messageCommandBuilder{c}
}

type userCommandBuilder struct {
	*commandBuilder
}

func (c userCommandBuilder) Option(option *objects.ApplicationCommandOption) UserCommandBuilder {
	c.commandBuilder.Option(option)
	return c
}

func (c userCommandBuilder) DefaultPermission() UserCommandBuilder {
	c.commandBuilder.DefaultPermission()
	return c
}

func (c userCommandBuilder) AllowedMentions(config *objects.AllowedMentions) UserCommandBuilder {
	c.commandBuilder.AllowedMentions(config)
	return c
}

func (c *commandBuilder) UserCommand() UserCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeUser)
	return userCommandBuilder{c}
}
