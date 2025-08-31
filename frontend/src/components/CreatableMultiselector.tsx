import React from "react";
import {
  MultiSelector,
  MultiSelectorContent,
  MultiSelectorInput,
  MultiSelectorItem,
  MultiSelectorList,
  MultiSelectorTrigger,
} from "@/components/ui/multi-select";

/** Narrow + normalize helpers */
function isString(v: unknown): v is string {
  return typeof v === "string";
}
const normalize = (s: string) => s.trim().replace(/\s+/g, " ");

/** Props */
export type CreatableMultiSelectorProps = {
  /** currently selected values (strings) */
  values: string[];
  /** selection change callback (will receive deduped, normalized strings) */
  onValuesChange: (vals: string[]) => void;

  /** available options shown in the list */
  options: string[];
  /** setter for available options (lets us append newly created ones) */
  setOptions: (opts: string[]) => void;

  placeholder?: string;
  className?: string;
};

/** Highly-optimized creatable wrapper */
export function CreatableMultiSelector({
  values,
  onValuesChange,
  options,
  setOptions,
  placeholder,
  className,
}: CreatableMultiSelectorProps) {
  const inputRef = React.useRef<HTMLInputElement | null>(null);
  const [rawInput, setRawInput] = React.useState("");

  // Dedup + normalize selected values (case-insensitive)
  const safeValues = React.useMemo(() => {
    const map = new Map<string, string>();
    for (const v of values ?? []) {
      if (isString(v)) {
        const n = normalize(v);
        if (n) map.set(n.toLowerCase(), n);
      }
    }
    return [...map.values()];
  }, [values]);

  // Dedup + normalize options (case-insensitive)
  const safeOptions = React.useMemo(() => {
    const map = new Map<string, string>();
    for (const o of options ?? []) {
      if (isString(o)) {
        const n = normalize(o);
        if (n) map.set(n.toLowerCase(), n);
      }
    }
    return [...map.values()];
  }, [options]);

  // Defer filter work to keep typing snappy
  const deferredQuery = React.useDeferredValue(rawInput);
  const filtered = React.useMemo(() => {
    const q = deferredQuery.toLowerCase();
    if (!q) return safeOptions;
    return safeOptions.filter((o) => o.toLowerCase().includes(q));
  }, [safeOptions, deferredQuery]);

  // Mark heavy updates as non-urgent so input stays responsive
  const [, startTransition] = React.useTransition();

  // Ensure dedup on any external change
  const handleValuesChange = (next: string[]) => {
    startTransition(() => {
      const map = new Map<string, string>();
      for (const v of next ?? []) {
        if (isString(v)) {
          const n = normalize(v);
          if (n) map.set(n.toLowerCase(), n);
        }
      }
      onValuesChange([...map.values()]);
    });
  };

  const add = (raw: string) => {
    const v = normalize(raw);
    if (!v) return;

    const existsOpt = safeOptions.some((o) => o.toLowerCase() === v.toLowerCase());
    const existsSel = safeValues.some((o) => o.toLowerCase() === v.toLowerCase());

    startTransition(() => {
      if (!existsOpt) setOptions([...safeOptions, v]);
      if (!existsSel) handleValuesChange([...safeValues, v]);
    });
    
    if (inputRef.current) inputRef.current.value = "";
    setRawInput("");
  };

  // NOTE: some MultiSelectorInput impls don't support onChange; use onInput
  const handleInput = (e: React.FormEvent<HTMLInputElement>) => {
    setRawInput((e.target as HTMLInputElement).value);
  };

  const tryCommit = () => {
    const v = inputRef.current?.value ?? rawInput;
    if (normalize(v)) add(v);
  };

  return (
    <MultiSelector
      values={safeValues}                // pass deduped selections
      onValuesChange={handleValuesChange}
      loop
      className={className}
    >
      <MultiSelectorTrigger>
        <MultiSelectorInput
          ref={inputRef}
          placeholder={placeholder ?? "Select or type…"}
          // keep it UNCONTROLLED for perf; track text via onInput
          onInput={handleInput}
          onKeyDown={(e) => {
            if (
              (e.key === "Enter" || e.key === "Tab" || e.key === ",") &&
              (inputRef.current?.value ?? "").trim()
            ) {
              e.preventDefault();
              tryCommit();
            }
          }}
          onBlur={tryCommit} // commit any leftover text when leaving the field
        />
      </MultiSelectorTrigger>

      <MultiSelectorContent>
        <MultiSelectorList>
          {filtered.map((opt) => (
            <MultiSelectorItem key={`opt-${opt.toLowerCase()}`} value={opt}>
              {opt}
            </MultiSelectorItem>
          ))}

          {/* Show "Create" affordance if input has text and it's not already an option */}
          {normalize(rawInput) &&
            !safeOptions.some((o) => o.toLowerCase() === rawInput.trim().toLowerCase()) && (
              <MultiSelectorItem
                key={`create-${normalize(rawInput).toLowerCase()}`}
                value={normalize(rawInput)}
                onClick={(e) => {
                  e.preventDefault();
                  add(rawInput);
                }}
              >
                Create “{normalize(rawInput)}”
              </MultiSelectorItem>
            )}
        </MultiSelectorList>
      </MultiSelectorContent>
    </MultiSelector>
  );
}

/* ------------------------------------------
   Example usage with React Hook Form (RHF)

   // In your form component:
   const form = useForm({
     defaultValues: { authors: [] }, // important
   });
   const [authorOptions, setAuthorOptions] = React.useState<string[]>(initialAuthors);

   <FormField
     control={form.control}
     name="authors"
     render={({ field }) => (
       <FormItem>
         <FormLabel>Authors</FormLabel>
         <FormControl>
           <CreatableMultiSelector
             values={field.value ?? []}
             onValuesChange={field.onChange}
             options={authorOptions}
             setOptions={setAuthorOptions}
             className="max-w-xs"
           />
         </FormControl>
         <FormMessage />
       </FormItem>
     )}
   />
------------------------------------------- */
