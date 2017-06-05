package collargo

/**
 * Public types
 */

// SendSignalFunc the function type used to send a signal type data
type SendSignalFunc func(Signal)

// SendDataFunc the function type used to send any type of data
type SendDataFunc func(interface{})

// FlowFunc the function converted from a flow
type FlowFunc func(data interface{}) (interface{}, error)

// Callback the callback function
type Callback func(err error, data interface{})

/**
 * Private types
 */
