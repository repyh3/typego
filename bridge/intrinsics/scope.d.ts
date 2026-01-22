declare namespace typego {
    /**
     * Defines a scope where deferred functions are executed LIFO upon exit.
     * Functions called with 'defer' inside this scope will be executed 
     * when the scope's callback returns or throws.
     */
    function scope<T>(fn: (defer: (cleanup: () => void) => void) => T): T;
}
