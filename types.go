package collargo

/**
 * Public types
 */

// SendSignalFunc the function type used to send a signal type data
type SendSignalFunc func(Signal)

// SendDataFunc the function type used to send any type of data
type SendDataFunc func(interface{})

// FlowFunc the function converted from a flow
type FlowFunc func(data interface{}, done Callback)

// Callback the nodejs style callback type
type Callback func(error, interface{})

/**
 * Private types
 */
