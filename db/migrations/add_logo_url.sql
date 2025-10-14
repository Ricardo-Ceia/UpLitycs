-- Migration: Add logo_url column to apps table
-- Run this if you have an existing database without the logo_url column

-- Add logo_url column if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name='apps' AND column_name='logo_url'
    ) THEN
        ALTER TABLE apps ADD COLUMN logo_url TEXT;
        RAISE NOTICE 'Added logo_url column to apps table';
    ELSE
        RAISE NOTICE 'logo_url column already exists';
    END IF;
END $$;
