#include <assert.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define AMBIENT_STACK_CAPACITY 1024
#define AMBIENT_PROGRAM_CAPACITY 1024
#define AMBIENT_EXECUTION_LIMIT 69

typedef int64_t Word;

typedef enum {
    ERROR_NOPE,
    ERROR_STACK_OVERFLOW,
    ERROR_STACK_UNDERFLOW,
    ERROR_ILLIGAL_INSTRUCTION,
    ERROR_ILLIGAL_INSTUCITION_ACCESS,
    ERROR_DIVIDE_BY_ZERO,
    ERROR_ILLIGAL_OPERAND,
} Error;

const char *error_to_string(Error error) {
    switch (error) {
    case ERROR_NOPE:
        return "nope";
    case ERROR_STACK_OVERFLOW:
        return "stack overflow";
    case ERROR_STACK_UNDERFLOW:
        return "stack underflow";
    case ERROR_ILLIGAL_INSTRUCTION:
        return "illigal instruction";
    case ERROR_ILLIGAL_INSTUCITION_ACCESS:
        return "illigal instruction access";
    case ERROR_DIVIDE_BY_ZERO:
        return "divide by zero";
    case ERROR_ILLIGAL_OPERAND:
        return "illigal operand";
    default:
        return "unknown error";
    }
}

typedef enum {
    INSTRUCTION_PUSH,
    INSTRUCTION_DUP,
    INSTRUCTION_PLUS,
    INSTRUCTION_MINUS,
    INSTRUCTION_MULTIPLY,
    INSTRUCTION_DIVISION,
    INSTRUCTION_JUMP,
    INSTRUCTION_JUMP_IF_TRUE,
    INSTRUCTION_EQUAL,
    INSTRUCTION_END,
    INSTRUCTION_PRINT_DEBUG,
} InstructionKind;

typedef struct {
    InstructionKind kind;
    Word operand;
} Instruction;

const char *instruction_to_string(InstructionKind kind) {
    switch (kind) {
    case INSTRUCTION_PUSH:
        return "push";
    case INSTRUCTION_DUP:
        return "dup";
    case INSTRUCTION_PLUS:
        return "plus";
    case INSTRUCTION_MINUS:
        return "minus";
    case INSTRUCTION_MULTIPLY:
        return "multiply";
    case INSTRUCTION_DIVISION:
        return "division";
    case INSTRUCTION_JUMP:
        return "jump";
    case INSTRUCTION_JUMP_IF_TRUE:
        return "jump if true";
    case INSTRUCTION_EQUAL:
        return "equal";
    case INSTRUCTION_END:
        return "end";
    case INSTRUCTION_PRINT_DEBUG:
        return "print debug";
    default:
        return "unknown instruction";
    }
}

typedef struct {
    Word stack[AMBIENT_STACK_CAPACITY];
    Word size;

    Instruction program[AMBIENT_PROGRAM_CAPACITY];
    Word current_instruction;
    Word program_size;

    bool end;
} Ambient;

void ambient_load_program_from_memory(Ambient *ambient,
                                      const Instruction *program, size_t size) {
    assert(size < AMBIENT_PROGRAM_CAPACITY);
    memcpy(ambient->program, program, sizeof(program[0]) * size);
    ambient->program_size = size;
}

#define MAKE_PUSH_INSTRUCTION(value)                                           \
    ((Instruction){.kind = INSTRUCTION_PUSH, .operand = (value)})

#define MAKE_JUMP_INSTRUCTION(address)                                         \
    ((Instruction){.kind = INSTRUCTION_JUMP, .operand = (address)})

#define MAKE_DUP_INSTRUCTION(address)                                          \
    ((Instruction){.kind = INSTRUCTION_DUP, .operand = (address)})

#define MAKE_END_INSTRUCTION() ((Instruction){.kind = INSTRUCTION_END})

#define MAKE_PLUS_INSTRUCTION() ((Instruction){.kind = INSTRUCTION_PLUS})

#define MAKE_MINUS_INSTRUCTION() ((Instruction){.kind = INSTRUCTION_MINUS})

#define MAKE_MULTIPLY_INSTRUCTION()                                            \
    ((Instruction){.kind = INSTRUCTION_MULTIPLY})

#define MAKE_DIVISION_INSTRUCTION()                                            \
    ((Instruction){.kind = INSTRUCTION_DIVISION})

#define CHECK_FOR_STACK_UNDERFLOW(glw)                                         \
    if ((glw)->size < 2) {                                                     \
        return ERROR_STACK_UNDERFLOW;                                          \
    }

Error ambient_execute_instruction(Ambient *ambient) {
    if (ambient->current_instruction < 0 ||
        ambient->current_instruction >= ambient->program_size) {
        return ERROR_ILLIGAL_INSTRUCTION;
    }

    Instruction instruction = ambient->program[ambient->current_instruction];

    switch (instruction.kind) {
    case INSTRUCTION_PUSH:
        if (ambient->size >= AMBIENT_STACK_CAPACITY) {
            return ERROR_STACK_OVERFLOW;
        }

        ambient->stack[ambient->size++] = instruction.operand;
        ambient->current_instruction++;

        break;
    case INSTRUCTION_PLUS:
        CHECK_FOR_STACK_UNDERFLOW(ambient);

        ambient->stack[ambient->size - 2] = ambient->stack[ambient->size - 2] +
                                            ambient->stack[ambient->size - 1];

        ambient->size -= 1;
        ambient->current_instruction++;

        break;
    case INSTRUCTION_MINUS:
        CHECK_FOR_STACK_UNDERFLOW(ambient);

        ambient->stack[ambient->size - 2] = ambient->stack[ambient->size - 2] -
                                            ambient->stack[ambient->size - 1];

        ambient->size -= 1;
        ambient->current_instruction++;

        break;
    case INSTRUCTION_MULTIPLY:
        CHECK_FOR_STACK_UNDERFLOW(ambient);

        ambient->stack[ambient->size - 2] = ambient->stack[ambient->size - 2] *
                                            ambient->stack[ambient->size - 1];

        ambient->size -= 1;
        ambient->current_instruction++;

        break;
    case INSTRUCTION_DIVISION:
        CHECK_FOR_STACK_UNDERFLOW(ambient);

        if (ambient->stack[ambient->size - 1] == 0) {
            return ERROR_DIVIDE_BY_ZERO;
        }

        ambient->stack[ambient->size - 2] = ambient->stack[ambient->size - 2] /
                                            ambient->stack[ambient->size - 1];

        ambient->size -= 1;
        ambient->current_instruction++;

        break;
    case INSTRUCTION_JUMP:
        ambient->current_instruction = instruction.operand;
        break;
    case INSTRUCTION_JUMP_IF_TRUE:
        CHECK_FOR_STACK_UNDERFLOW(ambient);

        if (ambient->stack[ambient->size - 1]) {
            ambient->size -= 1;
            ambient->current_instruction = instruction.operand;
        } else {
            ambient->current_instruction += 1;
        }

        break;
    case INSTRUCTION_EQUAL:
        CHECK_FOR_STACK_UNDERFLOW(ambient);

        ambient->stack[ambient->size - 2] = ambient->stack[ambient->size - 2] ==
                                            ambient->stack[ambient->size - 1];
        ambient->size -= 1;
        ambient->current_instruction += 1;

        break;
    case INSTRUCTION_PRINT_DEBUG:
        CHECK_FOR_STACK_UNDERFLOW(ambient);

        printf("%lld\n", ambient->stack[ambient->size - 1]);
        ambient->size -= 1;
        ambient->current_instruction += 1;

        break;
    case INSTRUCTION_DUP:
        if (ambient->size >= AMBIENT_STACK_CAPACITY) {
            return ERROR_STACK_OVERFLOW;
        }

        if (ambient->size - instruction.operand <= 0) {
            return ERROR_STACK_UNDERFLOW;
        }

        if (instruction.operand < 0) {
            return ERROR_ILLIGAL_OPERAND;
        }

        ambient->stack[ambient->size] =
            ambient->stack[ambient->size - 1 - instruction.operand];

        ambient->size += 1;
        ambient->current_instruction++;

        break;
    case INSTRUCTION_END:
        ambient->end = true;
        break;
    default:
        return ERROR_ILLIGAL_INSTRUCTION;
    }

    return ERROR_NOPE;
}

void ambient_dump(FILE *stream, const Ambient *ambient) {
    fprintf(stream, "stack:\n");
    if (ambient->size <= 0) {
        fprintf(stream, "    empty\n\n");
        return;
    }

    for (Word i = 0; i < ambient->size; i++) {
        fprintf(stream, "    %lld\n", ambient->stack[i]);
    }

    fprintf(stream, "\n");
}

Ambient ambient = {0};

Instruction program[] = {MAKE_PUSH_INSTRUCTION(0), MAKE_PUSH_INSTRUCTION(1),
                         MAKE_DUP_INSTRUCTION(1),  MAKE_DUP_INSTRUCTION(1),
                         MAKE_PLUS_INSTRUCTION(),  MAKE_JUMP_INSTRUCTION(2)};

#define ARRAY_SIZE(xs) (sizeof(xs) / sizeof((xs)[0]))

int main(void) {

    ambient_load_program_from_memory(&ambient, program, ARRAY_SIZE(program));

    for (int i = 0; i < 100 && !ambient.end; i++) {
        Error err = ambient_execute_instruction(&ambient);
        if (err != ERROR_NOPE) {
            fprintf(stderr, "error: %s\n", error_to_string(err));
            ambient_dump(stderr, &ambient);
            exit(1);
        }
    }

    ambient_dump(stdout, &ambient);

    return 0;
}
