package router

import (
	"fmt"
	"testing"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
)

func TestCommandRouterCtx_TargetMessage(t *testing.T) {
	tests := []struct {
		name string

		options map[string]interface{}
		expects *objects.Message
	}{
		{
			name: "no message",
		},
		{
			name: "message wrong type",
			options: map[string]interface{}{"/target": 1},
		},
		{
			name:    "message exists",
			options: map[string]interface{}{
				"/target": &ResolvableMessage{
					id:   "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Messages: map[objects.Snowflake]objects.Message{
								123: {
									Content: "hello world",
								},
							},
						},
					},
				},
			},
			expects: &objects.Message{Content: "hello world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.options == nil {
				tt.options = map[string]interface{}{}
			}
			ctx := CommandRouterCtx{Options: tt.options}
			assert.Equal(t, tt.expects, ctx.TargetMessage())
		})
	}
}

func TestCommandRouterCtx_TargetMember(t *testing.T) {
	tests := []struct {
		name string

		options map[string]interface{}
		expects *objects.GuildMember
	}{
		{
			name: "no member",
		},
		{
			name: "member wrong type",
			options: map[string]interface{}{"/target": 1},
		},
		{
			name:    "member exists",
			options: map[string]interface{}{
				"/target": &ResolvableUser{
					id:   "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Members: map[objects.Snowflake]objects.GuildMember{
								123: {
									Deaf: true,
								},
							},
						},
					},
				},
			},
			expects: &objects.GuildMember{Deaf: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.options == nil {
				tt.options = map[string]interface{}{}
			}
			ctx := CommandRouterCtx{Options: tt.options}
			assert.Equal(t, tt.expects, ctx.TargetMember())
		})
	}
}

func Test_messageTargetWrapper(t *testing.T) {
	tests := []struct {
		name string

		options map[string]interface{}
		expectsErr string
	}{
		{
			name: "successful call",
			options:  map[string]interface{}{
				"/target": &ResolvableMessage{
					id:   "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Messages: map[objects.Snowflake]objects.Message{
								123: {
									Content: "hello world",
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "invalid target",
			options: map[string]interface{}{},
			expectsErr: "wrong or no target specified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := &CommandRouterCtx{
				Options: tt.options,
			}

			var called bool
			cb := func(ctx *CommandRouterCtx, msg *objects.Message) error {
				if called {
					t.Fatal("function called twice")
				}
				called = true

				if mockCtx != ctx {
					t.Fatal("context different")
				}
				assert.Equal(t, &objects.Message{Content: "hello world"}, msg)

				return nil
			}

			f := messageTargetWrapper(cb)
			err := f(mockCtx)
			if tt.expectsErr == "" {
				if !called {
					t.Fatal("function wasn't called")
				}
				assert.NoError(t, err)
			} else {
				if called {
					t.Fatal("function was called")
				}
				assert.EqualError(t, err, tt.expectsErr)
			}
		})
	}
}

func Test_memberTargetWrapper(t *testing.T) {
	tests := []struct {
		name string

		options map[string]interface{}
		expectsErr string
	}{
		{
			name: "successful call",
			options:  map[string]interface{}{
				"/target": &ResolvableUser{
					id:   "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Members: map[objects.Snowflake]objects.GuildMember{
								123: {
									Nick: "hello world",
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "invalid target",
			options: map[string]interface{}{},
			expectsErr: "wrong or no target specified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := &CommandRouterCtx{
				Options: tt.options,
			}

			var called bool
			cb := func(ctx *CommandRouterCtx, member *objects.GuildMember) error {
				if called {
					t.Fatal("function called twice")
				}
				called = true

				if mockCtx != ctx {
					t.Fatal("context different")
				}
				assert.Equal(t, &objects.GuildMember{Nick: "hello world"}, member)

				return nil
			}

			f := memberTargetWrapper(cb)
			err := f(mockCtx)
			if tt.expectsErr == "" {
				if !called {
					t.Fatal("function wasn't called")
				}
				assert.NoError(t, err)
			} else {
				if called {
					t.Fatal("function was called")
				}
				assert.EqualError(t, err, tt.expectsErr)
			}
		})
	}
}

func TestCommandGroup_Use(t *testing.T) {
	tests := []struct {
		name string

		init func(c *CommandGroup)
		expectsLen int
	}{
		{
			name: "no middleware",
			expectsLen: 0,
		},
		{
			name: "add middleware",
			init: func(c *CommandGroup) {
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
			},
			expectsLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandGroup{Middleware: []MiddlewareFunc{}}
			if tt.init != nil {
				tt.init(c)
			}
			assert.Equal(t, tt.expectsLen, len(c.Middleware))
		})
	}
}

var dummyRootCommandGroup = &CommandGroup{}

var commandGroupTests = []struct {
	name string

	level uint
	groupName string
	description string
	defaultPermission bool

	expects *CommandGroup
	expectsErr string
}{
	{
		name:  "group nested too deep",
		level: 2,
		expectsErr: "sub-command group would be nested too deep",
	},
	{
		name: "root group",
		groupName: "abc",
		description: "def",
		defaultPermission: true,
		expects: &CommandGroup{
			level:             1,
			parent:            dummyRootCommandGroup,
			DefaultPermission: true,
			Description:       "def",
			Subcommands: map[string]interface{}{},
		},
	},
	{
		name: "sub-group",
		groupName: "abc",
		description: "def",
		defaultPermission: true,
		level: 1,
		expects: &CommandGroup{
			level:             2,
			parent:            dummyRootCommandGroup,
			DefaultPermission: true,
			Description:       "def",
			Subcommands: map[string]interface{}{},
		},
	},
}

func TestCommandGroup_NewCommandGroup(t *testing.T) {
	for _, tt := range commandGroupTests {
		t.Run(tt.name, func(t *testing.T) {
			dummyRootCommandGroup.Subcommands, dummyRootCommandGroup.level = map[string]interface{}{}, tt.level
			group, err := dummyRootCommandGroup.NewCommandGroup(tt.groupName, tt.description, tt.defaultPermission)
			if tt.expectsErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectsErr)
				return
			}
			assert.Equal(t, tt.expects, group)
		})
	}
}

type mustNewCommandGroup interface {
	MustNewCommandGroup(name, description string, defaultPermission bool) *CommandGroup
}

func unpanicCommandGroup(x mustNewCommandGroup, name, description string, default_ bool) (group *CommandGroup, returnedErr string) {
	defer func() {
		if r := recover(); r != nil {
			returnedErr = fmt.Sprint(r)
		}
	}()
	group = x.MustNewCommandGroup(name, description, default_)
	return
}

func TestCommandGroup_MustNewCommandGroup(t *testing.T) {
	for _, tt := range commandGroupTests {
		t.Run(tt.name, func(t *testing.T) {
			dummyRootCommandGroup.Subcommands, dummyRootCommandGroup.level = map[string]interface{}{}, tt.level
			group, errResult := unpanicCommandGroup(dummyRootCommandGroup, tt.groupName, tt.description, tt.defaultPermission)
			assert.Equal(t, errResult, tt.expectsErr)
			assert.Equal(t, tt.expects, group)
		})
	}
}

func TestCommandRouter_Use(t *testing.T) {
	tests := []struct {
		name string

		init func(c *CommandRouter)
		expectsLen int
	}{
		{
			name: "no middleware",
			expectsLen: 0,
		},
		{
			name: "add middleware",
			init: func(c *CommandRouter) {
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
			},
			expectsLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CommandRouter{middleware: []MiddlewareFunc{}}
			if tt.init != nil {
				tt.init(r)
			}
			assert.Equal(t, tt.expectsLen, len(r.middleware))
		})
	}
}

func TestCommandRouter_NewCommandGroup(t *testing.T) {
	r := &CommandRouter{}
	group, err := r.NewCommandGroup("abc", "def", true)
	assert.NoError(t, err)
	assert.Equal(t, &CommandGroup{
		level:             1,
		DefaultPermission: true,
		Description:       "def",
		Subcommands:       map[string]interface{}{},
	}, group)
}

func TestCommandRouter_MustNewCommandGroup(t *testing.T) {
	r := &CommandRouter{}
	group, errResult := unpanicCommandGroup(r, "abc", "def", true)
	assert.Equal(t, "", errResult)
	assert.Equal(t, &CommandGroup{
		level:             1,
		DefaultPermission: true,
		Description:       "def",
		Subcommands:       map[string]interface{}{},
	}, group)
}

func TestCommandRouter_NewCommandBuilder(t *testing.T) {
	r := &CommandRouter{}
	builder := r.NewCommandBuilder("abc")
	assert.NotNil(t, r.roots.Subcommands)
	assert.Equal(t, &commandBuilder{
		map_: r.roots.Subcommands,
		cmd:  Command{Name: "abc"},
	}, builder)
}

func TestCommandGroup_execute(t *testing.T) {
	// TODO
}

func TestCommandRouter_build(t *testing.T) {
	// TODO
}

func Test_getOptions(t *testing.T) {
	// TODO
}

func TestCommandRouter_FormulateDiscordCommands(t *testing.T) {
	// TODO
}

func TestCommandRouterCtx_Bind(t *testing.T) {
	// TODO
}