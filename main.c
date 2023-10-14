#include <assert.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define GLOW_STACK_CAPACITY 1024
#define GLOW_PROGRAM_CAPACITY 1024
#define GLOW_EXECUTION_LIMIT 69

typedef int64_t Word;

typedef enum {
    ERROR_NOPE,
    ERROR_STACK_OVERFLOW,
    ERROR_STACK_UNDERFLOW,
    ERROR_ILLIGAL_INSTRUCTION,
    ERROR_ILLIGAL_INSTUCITION_ACCESS,
    ERROR_DIVIDE_BY_ZERO,
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
    default:
        return "unknown error";
    }
}

typedef enum {
    INSTRUCTION_PUSH,
    INSTRUCTION_PLUS,
    INSTRUCTION_MINUS,
    INSTRUCTION_MULTIPLY,
    INSTRUCTION_DIVISION,
    INSTRUCTION_JUMP,
    INSTRUCTION_JUMP_IF_TRUE,
    INSTRUCTION_EQUAL,
    INSTRUCTION_END,
} InstructionKind;

typedef struct {
    InstructionKind kind;
    Word operand;
} Instruction;

const char *instruction_to_string(InstructionKind kind) {
    switch (kind) {
    case INSTRUCTION_PUSH:
        return "push";
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
    default:
        return "unknown instruction";
    }
}

typedef struct {
    Word stack[GLOW_STACK_CAPACITY];
    Word size;

    Instruction program[GLOW_PROGRAM_CAPACITY];
    Word current_instruction;
    Word program_size;

    bool end;
} Glow;

void glow_load_program_from_memory(Glow *glow, const Instruction *program,
                                   size_t size) {
    assert(size < GLOW_PROGRAM_CAPACITY);
    memcpy(glow->program, program, sizeof(program[0]) * size);
    glow->program_size = size;
}

#define MAKE_PUSH_INSTRUCTION(value)                                           \
    ((Instruction){.kind = INSTRUCTION_PUSH, .operand = (value)})

#define MAKE_JUMP_INSTRUCTION(address)                                         \
    ((Instruction){.kind = INSTRUCTION_JUMP, .operand = (address)})

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

Error glow_execute_instruction(Glow *glow) {
    if (glow->current_instruction < 0 ||
        glow->current_instruction >= glow->program_size) {
        return ERROR_ILLIGAL_INSTRUCTION;
    }

    Instruction instruction = glow->program[glow->current_instruction];

    switch (instruction.kind) {
    case INSTRUCTION_PUSH:
        if (glow->size >= GLOW_STACK_CAPACITY) {
            return ERROR_STACK_OVERFLOW;
        }

        glow->stack[glow->size++] = instruction.operand;
        glow->current_instruction++;

        break;
    case INSTRUCTION_PLUS:
        CHECK_FOR_STACK_UNDERFLOW(glow);

        glow->stack[glow->size - 2] =
            glow->stack[glow->size - 2] + glow->stack[glow->size - 1];

        glow->size -= 1;
        glow->current_instruction++;

        break;
    case INSTRUCTION_MINUS:
        CHECK_FOR_STACK_UNDERFLOW(glow);

        glow->stack[glow->size - 2] =
            glow->stack[glow->size - 2] - glow->stack[glow->size - 1];

        glow->size -= 1;
        glow->current_instruction++;

        break;
    case INSTRUCTION_MULTIPLY:
        CHECK_FOR_STACK_UNDERFLOW(glow);

        glow->stack[glow->size - 2] =
            glow->stack[glow->size - 2] * glow->stack[glow->size - 1];

        glow->size -= 1;
        glow->current_instruction++;

        break;
    case INSTRUCTION_DIVISION:
        CHECK_FOR_STACK_UNDERFLOW(glow);

        if (glow->stack[glow->size - 1] == 0) {
            return ERROR_DIVIDE_BY_ZERO;
        }

        glow->stack[glow->size - 2] =
            glow->stack[glow->size - 2] / glow->stack[glow->size - 1];

        glow->size -= 1;
        glow->current_instruction++;

        break;
    case INSTRUCTION_JUMP:
        glow->current_instruction = instruction.operand;
        break;
    case INSTRUCTION_JUMP_IF_TRUE:

        break;
    case INSTRUCTION_EQUAL:

        break;
    case INSTRUCTION_END:
        glow->end = true;
        break;
    default:
        return ERROR_ILLIGAL_INSTRUCTION;
    }

    return ERROR_NOPE;
}

bool _check_for_stackoverflow(Glow *glow) {
    if (glow->size >= GLOW_STACK_CAPACITY) {
        return true;
    }

    return false;
}

void glow_dump(FILE *stream, const Glow *glow) {
    fprintf(stream, "stack:\n");
    if (glow->size <= 0) {
        fprintf(stream, "    empty\n\n");
        return;
    }

    for (Word i = 0; i < glow->size; i++) {
        fprintf(stream, "    %lld\n", glow->stack[i]);
    }

    fprintf(stream, "\n");
}

Glow glow = {0};

Instruction program[] = {
    MAKE_PUSH_INSTRUCTION(0),
    MAKE_PUSH_INSTRUCTION(1),
    MAKE_PLUS_INSTRUCTION(),
    MAKE_JUMP_INSTRUCTION(1),
};

#define ARRAY_SIZE(xs) (sizeof(xs) / sizeof((xs)[0]))

int main(void) {

    glow_load_program_from_memory(&glow, program, ARRAY_SIZE(program));

    for (int i = 0; i < GLOW_EXECUTION_LIMIT && !glow.end; i++) {
        Error err = glow_execute_instruction(&glow);
        if (err != ERROR_NOPE) {
            fprintf(stderr, "error: %s\n", error_to_string(err));
            glow_dump(stderr, &glow);
            exit(1);
        }
    }

    glow_dump(stdout, &glow);

    return 0;
}
