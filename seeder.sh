#!/usr/bin/env bash
set -e

API="${API:-http://localhost:8080}"
DB_PATH="${DB_PATH:-backend/engager.db}"

echo "=========================================================="
echo "1. Cleaning previous records from SQLite..."
echo "=========================================================="

if [ -f "$DB_PATH" ]; then
  sqlite3 "$DB_PATH" << 'EOF'
DELETE FROM personalized_activities WHERE 1=1;
DELETE FROM activity_attempts WHERE 1=1;
DELETE FROM activities WHERE 1=1;
VACUUM;
EOF
  echo "Cleaned existing activities, attempts, and personalized records from $DB_PATH."
else
  echo "Database file '$DB_PATH' not found locally; continuing with API creation..."
fi

echo
echo "=========================================================="
echo "2. Seeding curated activities via $API/api/activities..."
echo "=========================================================="

create_activity() {
  local payload="$1"
  local title
  title="$(echo "$payload" | python3 -c 'import sys, json; print(json.load(sys.stdin).get("title", "Untitled"))')"
  
  response=$(curl -sS -w "\n%{http_code}" -X POST "$API/api/activities" \
    -H "Content-Type: application/json" \
    -d "$payload")

  status_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | head -n -1)

  if [ "$status_code" -ge 200 ] && [ "$status_code" -lt 300 ]; then
    echo "  [OK] ($status_code) Created: $title"
  else
    echo "  [ERR] ($status_code) Failed to create '$title': $body"
    exit 1
  fi
}

# ------------------------------------------------------------------------------
# TRACK 1: Computer Science & Programming
# ------------------------------------------------------------------------------

create_activity '{
  "title": "Recursion Basics: Factorial",
  "description": "Understand recursive base cases and reduction steps.",
  "content": "Write a recursive function to compute the factorial of a positive integer n. Clearly state your base case condition and your recursive recurrence step.",
  "correct_answer": "Base case is n <= 1 returning 1; recursive step is n * factorial(n - 1)",
  "topic": "computer-science",
  "difficulty": 1,
  "estimated_duration_minutes": 3
}'

create_activity '{
  "title": "Binary Tree Path Traversal",
  "description": "Traverse nested hierarchical data structures.",
  "content": "Given a binary tree where each node holds an integer value, explain how to recursively determine the maximum root-to-leaf path sum.",
  "correct_answer": "Recursively compute node.val + max(maxPath(left), maxPath(right)), returning 0 when a child node is null.",
  "topic": "computer-science",
  "difficulty": 2,
  "estimated_duration_minutes": 5
}'

create_activity '{
  "title": "N-Queens Recursive Backtracking",
  "description": "State validation and search pruning under constraints.",
  "content": "In solving the N-Queens problem on an N×N board using backtracking, how can you validate whether a candidate column and both diagonals are threatened in O(1) time?",
  "correct_answer": "Maintain lookup hash sets or bitmasks tracking column indices, (row - col) for major diagonals, and (row + col) for minor diagonals.",
  "topic": "computer-science",
  "difficulty": 4,
  "estimated_duration_minutes": 8
}'

# ------------------------------------------------------------------------------
# TRACK 2: Mathematics & Problem Solving
# ------------------------------------------------------------------------------

create_activity '{
  "title": "Linear Equation Isolation",
  "description": "Isolate single variables using inverse operations.",
  "content": "Solve for x: 3x + 15 = 42. Show the inverse arithmetic operations needed to isolate x.",
  "correct_answer": "Subtract 15 to get 3x = 27, then divide by 3 to get x = 9.",
  "topic": "mathematics",
  "difficulty": 1,
  "estimated_duration_minutes": 2
}'

create_activity '{
  "title": "Systems of Linear Equations",
  "description": "Solve multi-variable systems via substitution or elimination.",
  "content": "Solve the system: 2x + y = 14 and 4x - 2y = 8. What is the coordinate solution (x, y)?",
  "correct_answer": "(4, 6)",
  "topic": "mathematics",
  "difficulty": 2,
  "estimated_duration_minutes": 4
}'

create_activity '{
  "title": "Mental Math: The Difference of Squares Trick",
  "description": "Apply algebraic identities to compute large values mentally.",
  "content": "Using the algebraic identity (a - b)(a + b) = a^2 - b^2, calculate 999 × 1001 without using long multiplication.",
  "correct_answer": "(1000 - 1)(1000 + 1) = 1,000,000 - 1 = 999,999",
  "topic": "mathematics",
  "difficulty": 3,
  "estimated_duration_minutes": 3
}'

create_activity '{
  "title": "Matrix Determinants & Linear Invertibility",
  "description": "Geometric properties of 2x2 matrix transformations.",
  "content": "Given matrix A = [[3, 2], [1, 4]], calculate the determinant det(A). What does a non-zero determinant signify regarding the invertibility of the transformation?",
  "correct_answer": "det(A) = (3*4) - (2*1) = 10. A non-zero determinant indicates the matrix is non-singular and strictly invertible.",
  "topic": "mathematics",
  "difficulty": 4,
  "estimated_duration_minutes": 6
}'

# ------------------------------------------------------------------------------
# TRACK 3: Logic, Deductions & Pop Culture
# ------------------------------------------------------------------------------

create_activity '{
  "title": "Emoji Movie Detective",
  "description": "Decode famous cinema titles from visual clues.",
  "content": "Decode this classic animated movie title represented by these emojis: 🦁👑🌅",
  "correct_answer": "The Lion King",
  "topic": "pop-culture",
  "difficulty": 1,
  "estimated_duration_minutes": 2
}'

create_activity '{
  "title": "The Mystery Number",
  "description": "Deduce numbers using digit constraints.",
  "content": "I am a two-digit number. My digits sum to 9. My tens digit is twice my ones digit. What number am I?",
  "correct_answer": "63",
  "topic": "logic",
  "difficulty": 2,
  "estimated_duration_minutes": 3
}'

create_activity '{
  "title": "The Three Light Switches Paradox",
  "description": "Solve constraints using state changes and heat.",
  "content": "You are outside a closed room with 3 switches controlling 1 incandescent light bulb inside. You may enter the room only once. How do you definitively determine which switch controls the bulb?",
  "correct_answer": "Turn switch 1 ON for 5 minutes, turn it OFF, turn switch 2 ON, and enter. The lighted bulb corresponds to switch 2, the warm unlit bulb to switch 1, and the cold unlit bulb to switch 3.",
  "topic": "logic",
  "difficulty": 4,
  "estimated_duration_minutes": 5
}'

# ------------------------------------------------------------------------------
# TRACK 4: Creativity & Systems Thinking
# ------------------------------------------------------------------------------

create_activity '{
  "title": "One-Minute Invention",
  "description": "Rapid creative problem solving.",
  "content": "You have 60 seconds to invent a device that solves an annoying everyday morning routine issue. Give it a name and describe how it works in two sentences.",
  "correct_answer": "Any coherent and creative invention description.",
  "topic": "creativity",
  "difficulty": 1,
  "estimated_duration_minutes": 2
}'

create_activity '{
  "title": "Perspective Shift Narrative",
  "description": "Re-examine established narratives from alternative viewpoints.",
  "content": "Take the classic fairy tale of The Three Little Pigs. Write a 3-sentence summary of the conflict told entirely from the perspective of the Big Bad Wolf.",
  "correct_answer": "Any coherent creative response expressing the antagonist perspective.",
  "topic": "creativity",
  "difficulty": 2,
  "estimated_duration_minutes": 4
}'

create_activity '{
  "title": "Island Civilization Governance",
  "description": "Balancing interdependent systems and institutional trade-offs.",
  "content": "You are founding a new island community. Define: 1) Its primary energy resource, 2) Its supreme law, and 3) Its core societal virtue. Explain in 2 sentences how these three decisions support one another.",
  "correct_answer": "Any coherent, reasoned systems-thinking response.",
  "topic": "systems-thinking",
  "difficulty": 3,
  "estimated_duration_minutes": 6
}'

echo
echo "=========================================================="
echo "Seeding completed successfully! All activities are ready."
echo "=========================================================="
