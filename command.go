package router

import (
	"container/list"
	"errors"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
	"strconv"
)

// Command is used to define a Discord (sub-)command. DO NOT MAKE YOURSELF! USE CommandGroup.NewCommandBuilder OR CommandRouter.NewCommandBuilder!
type Command struct {
	// Defines the command type.
	commandType int

	// Defines the parent.
	parent *CommandGroup

	// Description is the description for the command.
	Description string `json:"description"`

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions"`

	// DefaultPermission is used to define if the permissions of the command are default.
	DefaultPermission bool `json:"default_permission"`

	// Options defines the options which are required for a command.
	Options []*objects.ApplicationCommandOption `json:"options"`

	// Function is used to define the command being called.
	Function func(*CommandRouterCtx) error `json:"-"`
}

// Finds the option.
func findOption(name string, options []*objects.ApplicationCommandOption) *objects.ApplicationCommandOption {
	if options == nil {
		return nil
	}
	for _, v := range options {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// NonExistentOption is thrown when an option is provided in an interaction that doesn't exist in the command.
var NonExistentOption = errors.New("interaction option doesn't exist on command")

// MismatchedOption is thrown when the option types mismatch.
var MismatchedOption = errors.New("mismatched interaction option")

// Defines the options for command execution.
type commandExecutionOptions struct {
	restClient       *rest.Client
	exceptionHandler func(error) *objects.InteractionResponse
	allowedMentions  *objects.AllowedMentions
	interaction      *objects.Interaction
	data             *objects.ApplicationCommandInteractionData
	options          []*objects.ApplicationCommandInteractionDataOption
}

// Execute the command.
func (c *Command) execute(opts commandExecutionOptions, middlewareList *list.List) *objects.InteractionResponse {
	// Process the options.
	mappedOptions := map[string]interface{}{}

	if opts.data.TargetID != 0 {
		// Add a special case for "/target". The slash is there as a keyword.
		if _, ok := opts.data.Resolved.Messages[opts.data.TargetID]; ok {
			mappedOptions["/target"] = &ResolvableMessage{
				id:   strconv.FormatUint(uint64(opts.data.TargetID), 10),
				data: opts.data,
			}
		} else if _, ok = opts.data.Resolved.Users[opts.data.TargetID]; ok {
			mappedOptions["/target"] = &ResolvableUser{
				id:   strconv.FormatUint(uint64(opts.data.TargetID), 10),
				data: opts.data,
			}
		}
	} else {
		for _, v := range opts.options {
			// Find the option.
			option := findOption(v.Name, c.Options)
			if option == nil {
				// Option was provided that was not in this command.
				return opts.exceptionHandler(NonExistentOption)
			}

			// Check the option and result match types.
			if objects.ApplicationCommandOptionType(v.Type) != option.OptionType {
				return opts.exceptionHandler(MismatchedOption)
			}

			// Check what the type is.
			switch option.OptionType {
			case objects.TypeChannel:
				mappedOptions[option.Name] = &ResolvableChannel{
					id:   v.Value.(string),
					data: opts.data,
				}
			case objects.TypeRole:
				mappedOptions[option.Name] = &ResolvableRole{
					id:   v.Value.(string),
					data: opts.data,
				}
			case objects.TypeUser:
				mappedOptions[option.Name] = &ResolvableUser{
					id:   v.Value.(string),
					data: opts.data,
				}
			case objects.TypeString:
				mappedOptions[option.Name] = v.Value.(string)
			case objects.TypeInteger:
				mappedOptions[option.Name] = int(v.Value.(float64))
			case objects.TypeBoolean:
				mappedOptions[option.Name] = v.Value.(bool)
			case objects.TypeMentionable:
				mappedOptions[option.Name] = &ResolvableMentionable{
					id:   v.Value.(string),
					data: opts.data,
				}
			case objects.TypeDouble:
				mappedOptions[option.Name] = v.Value.(float64)
			}
		}
	}

	// Get the allowed mentions configuration.
	if c.AllowedMentions != nil {
		opts.allowedMentions = c.AllowedMentions
	}

	// Attempt to catch errors from here.
	defer func() {
		if errGeneric := recover(); errGeneric != nil {
			// Shouldn't try and return from defer.
			opts.exceptionHandler(ungenericError(errGeneric))
		}
	}()

	// Create the context.
	rctx := &CommandRouterCtx{
		globalAllowedMentions: opts.allowedMentions,
		errorHandler:          opts.exceptionHandler,
		Interaction:           opts.interaction,
		Command:               c,
		Options:               mappedOptions,
		RESTClient:            opts.restClient,
	}

	// Run the command.
	handler := c.Function
	if handler == nil {
		// don't nil crash
		handler = func(ctx *CommandRouterCtx) error { return nil }
	}
	if middlewareList.Len() == 0 {
		// Just call the command function.
		if err := handler(rctx); err != nil {
			return opts.exceptionHandler(err)
		}
	} else {
		// Wrap the end command function in a middleware function and push it.
		middlewareWrapper := func(ctx MiddlewareCtx) error {
			return handler(ctx.CommandRouterCtx)
		}
		middlewareList.PushBack(middlewareWrapper)

		// Create the middleware context.
		mctx := MiddlewareCtx{CommandRouterCtx: rctx, middlewareList: middlewareList}

		// Call the middleware chain.
		if err := mctx.Next(); err != nil {
			return opts.exceptionHandler(err)
		}
	}
	return rctx.buildResponse(false, opts.exceptionHandler, opts.allowedMentions)
}

// Groups is used to get the command groups that this belongs to.
func (c *Command) Groups() []*CommandGroup {
	s := []*CommandGroup{}
	for x := c.parent; x != nil; x = x.parent {
		s = append([]*CommandGroup{x}, s...)
	}
	return s
}
