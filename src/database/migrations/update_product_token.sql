-- Add created_by_id and is_active columns
ALTER TABLE product_tokens 
ADD COLUMN IF NOT EXISTS created_by_id UUID DEFAULT NULL,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;

-- Fix existing records with invalid user_id references
UPDATE product_tokens 
SET user_id = NULL 
WHERE user_id IS NOT NULL AND user_id NOT IN (SELECT id FROM users);

-- Add foreign key constraint for user_id
ALTER TABLE product_tokens 
ADD CONSTRAINT IF NOT EXISTS fk_product_tokens_user 
FOREIGN KEY (user_id) 
REFERENCES users(id) 
ON DELETE SET NULL;

-- Add foreign key constraint for created_by_id
ALTER TABLE product_tokens 
ADD CONSTRAINT IF NOT EXISTS fk_product_tokens_created_by 
FOREIGN KEY (created_by_id) 
REFERENCES users(id) 
ON DELETE SET NULL; 