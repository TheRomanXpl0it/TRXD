"use client"
import {
  useState,
  useMemo
} from "react"
import {
  toast
} from "sonner"
import {
  useForm,
  useFieldArray
} from "react-hook-form"
import {
  zodResolver
} from "@hookform/resolvers/zod"
import * as z from "zod"
import {
  Button
} from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import {
  Input
} from "@/components/ui/input"
import {
  Textarea
} from "@/components/ui/textarea"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select"
import {
  TagsInput
} from "@/components/ui/tags-input"
import {
  Switch
} from "@/components/ui/switch"
import {
  CloudUpload,
  Paperclip,
  Container,
  MemoryStick,
  Cpu
} from "lucide-react"
import {
  FileInput,
  FileUploader,
  FileUploaderContent,
  FileUploaderItem
} from "@/components/ui/file-upload"
import {
  MultiSelector,
  MultiSelectorContent,
  MultiSelectorInput,
  MultiSelectorItem,
  MultiSelectorList,
  MultiSelectorTrigger
} from "@/components/ui/multi-select"
import { AuthProps } from "@/context/AuthProvider"
import { Challenge as ChallengeType } from "@/context/ChallengeProvider"
import type { DockerConfig } from "@/context/ChallengeProvider"
import { oneDark } from "@codemirror/theme-one-dark";


// New: CodeMirror for syntax highlighting
import CodeMirror from "@uiw/react-codemirror"
import { yaml } from "@codemirror/lang-yaml"


const formSchema = z.object({
  title: z.string().min(1),
  flags: z.array(
    z.object({
      flag: z.string().min(1),
      regex: z.boolean(),
    })
  ).min(1, "At least one flag required"),
  description: z.string(),
  difficulty: z.string(),
  tags: z.array(z.string()).nonempty("Please at least one item"),
  hidden: z.boolean(),
  files: z.string(),
  authors: z.array(z.string()),
  instanced: z.boolean(),
  docker: z.object({
    image: z.string().min(1, "Docker image is required"),
    compose: z.string().optional().default(""),
    hashDomain: z.boolean().optional().default(false),
    lifetime: z.number().int().nonnegative().default(3600),
    envs: z.string().optional().default(""),
    maxMemory: z.number().nonnegative().optional().default(0),
    maxCPU: z.number().nonnegative().optional().default(0),
  }).optional(),
}).superRefine((val, ctx) => {
  if (val.instanced && !val.docker) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ["docker"],
      message: "Docker configuration is required when Instanced is enabled.",
    });
  }
});

function displayAuthors(authors: string[]){
  return (
    authors.map((author, index) => (
      <MultiSelectorItem key={index} value={author}>{author}</MultiSelectorItem>
    ))
  )
}

export function ChallengeForm( { challenge, auth } : { 
  challenge?: ChallengeType,
  auth: AuthProps
}) {

  let defaultAuthors: string[] | undefined = [auth.username];
  let defaultTitle: string | undefined = "";
  let defaultDescription: string | undefined = "";
  let defaultDifficulty: string | undefined = "Easy";
  let defaultTags: string[] | undefined = ["Easy"];
  let defaultHidden: boolean | undefined = true;
  let defaultFlags: { flag: string; regex: boolean }[] = [{ flag: "", regex: false }];
  let defaultInstanced: boolean | undefined = false;
  let dockerConfig: DockerConfig | undefined = undefined;
  let authors: string[] = [];

  if ( challenge ){
    defaultAuthors = challenge.authors ?? [auth.username];
    defaultDescription = challenge.description ?? "";
    defaultDifficulty = challenge.difficulty ?? "Easy";
    defaultTags = challenge.tags ?? ["Easy"];
    defaultHidden = challenge.hidden ?? true;
    defaultTitle = challenge.title ?? "";
    defaultInstanced = challenge.instanced ?? false;
    dockerConfig = challenge.docker ?? undefined;
    
    if (challenge.flags?.length) {
      defaultFlags = challenge.flags.map(flag =>
        typeof flag === "string" ? { flag: flag, regex: false } : flag
      );
    }
  }

  const [files, setFiles] = useState < File[] | null > (null);
  
  const dropZoneConfig = {
    maxFiles: 5,
    maxSize: 1024 * 1024 * 4,
    multiple: true,
  };

  const form = useForm < z.infer < typeof formSchema >> ({
    resolver: zodResolver(formSchema),
    defaultValues: {
      "title" : defaultTitle,
      "flags": defaultFlags,
      "description" : defaultDescription,
      "difficulty" : defaultDifficulty,
      "tags": defaultTags,
      "authors": defaultAuthors,
      "hidden": defaultHidden,
      "instanced": defaultInstanced,
      "docker": dockerConfig
        ? {
            image: dockerConfig.image ?? "",
            compose: (dockerConfig as any).compose ?? "",
            hashDomain: (dockerConfig as any).hashDomain ?? false,
            lifetime: (dockerConfig as any).lifetime ?? 3600,
            envs: (dockerConfig as any).envs ?? "",
            maxMemory: (dockerConfig as any).maxMemory ?? 0,
            maxCPU: (dockerConfig as any).maxCPU ?? 0,
          }
        : undefined,
    },
  })
  
  const { control, watch, setValue } = form;
  const { fields: flagFields, append, remove } = useFieldArray({
    control,
    name: "flags",
  });

  const instanced = watch("instanced");
  const yamlExtensions = useMemo(() => [yaml()], []);

  function onSubmit(values: z.infer < typeof formSchema > ) {
    try {
      const payload = {
        ...values,
        docker: values.instanced ? values.docker : undefined,
      }
      console.log(payload);
      toast(
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(payload, null, 2)}</code>
        </pre>
      );
    } catch (error) {
      console.error("Form submission error", error);
      toast.error("Failed to submit the form. Please try again.");
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 mx-auto py-10">
        
      <FormField
        control={form.control}
        name="title"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Title</FormLabel>
            <FormControl>
              <Input 
              placeholder=""
              
              type="text"
              {...field} />
            </FormControl>
            <FormDescription>The title of the challenge</FormDescription>
            <FormMessage />
          </FormItem>
        )}
        />

        <FormLabel>Flags</FormLabel>
        {flagFields.map((field, index) => (
        <div key={field.id} className="space-y-2 border p-4 rounded-md">
          <FormField
            control={form.control}
            name={`flags.${index}.flag`}
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-sm text-muted-foreground">Flag {index + 1}</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="TRX{...}" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name={`flags.${index}.regex`}
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
                  <FormLabel>Regex</FormLabel>
                  <FormDescription>Enable if this flag should be matched as a regular expression.</FormDescription>
                </div>
                <FormControl>
                  <Switch
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
              </FormItem>
            )}
          />
          {flagFields.length > 1 && (
            <Button type="button" variant="destructive" onClick={() => remove(index)}>
              Remove
            </Button>
          )}
        </div>
        ))}
        <Button type="button" onClick={() => append({ flag: "", regex: false })}>
        + Add Flag
        </Button>

        <FormField
        control={form.control}
        name="description"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Description</FormLabel>
            <FormControl>
              <Textarea
                placeholder=""
                className="resize-none"
                {...field}
              />
            </FormControl>
            <FormDescription>The description of the challenge</FormDescription>
            <FormMessage />
          </FormItem>
        )}
        />

        <FormField
        control={form.control}
        name="difficulty"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Difficulty</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select difficulty" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectItem value="Easy">Easy</SelectItem>
                <SelectItem value="Medium">Medium</SelectItem>
                <SelectItem value="Hard">Hard</SelectItem>
                <SelectItem value="Insane">Insane</SelectItem>
              </SelectContent>
            </Select>
              <FormDescription>Difficulty of the challenge</FormDescription>
            <FormMessage />
          </FormItem>
        )}
        />

        <FormField
        control={form.control}
        name="tags"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Tags</FormLabel>
            <FormControl>
              <TagsInput
                value={field.value}
                onValueChange={field.onChange}
                placeholder="Enter your tags"
              />
            </FormControl>
            <FormDescription>Add tags to the challenge.</FormDescription>
            <FormMessage />
          </FormItem>
        )}
        />

        <FormField
            control={form.control}
            name="hidden"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
                  <FormLabel>Hidden</FormLabel>
                  <FormDescription>Challenge is hidden, you can edit it later on inside the challenge card.</FormDescription>
                </div>
                <FormControl>
                  <Switch
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="files"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Challenge files</FormLabel>
                <FormControl>
                  <FileUploader
                    value={files}
                    onValueChange={setFiles}
                    dropzoneOptions={dropZoneConfig}
                    className="relative bg-background rounded-lg p-2"
                  >
                    <FileInput
                      id="fileInput"
                      className="outline-dashed outline-1 outline-slate-500"
                    >
                      <div className="flex items-center justify-center flex-col p-8 w-full ">
                        <CloudUpload className='text-gray-500 w-10 h-10' />
                        <p className="mb-1 text-sm text-gray-500 dark:text-gray-400">
                          <span className="font-semibold">Click to upload</span>
                          &nbsp; or drag and drop
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                          (.zip)
                        </p>
                      </div>
                    </FileInput>
                    <FileUploaderContent>
                      {files &&
                        files.length > 0 &&
                        files.map((file, i) => (
                          <FileUploaderItem key={i} index={i}>
                            <Paperclip className="h-4 w-4 stroke-current" />
                            <span>{file.name}</span>
                          </FileUploaderItem>
                        ))}
                    </FileUploaderContent>
                  </FileUploader>
                </FormControl>
                <FormDescription>Select a file to upload.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="authors"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Authors</FormLabel>
                <FormControl>
                  <MultiSelector
                    values={field.value}
                    onValuesChange={field.onChange}
                    loop
                    className="max-w-xs"
                  >
                    <MultiSelectorTrigger>
                      <MultiSelectorInput placeholder="Select Authors" />
                    </MultiSelectorTrigger>
                    <MultiSelectorContent>
                    <MultiSelectorList>
                      { displayAuthors(authors) }
                    </MultiSelectorList>
                    </MultiSelectorContent>
                  </MultiSelector>
                </FormControl>
                <FormDescription>Select multiple options.</FormDescription>
                <FormMessage />
              </FormItem>
                  )}
              />

                <FormField
                  control={form.control}
                  name="instanced"
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                      <div className="space-y-0.5">
                        <FormLabel>Instanced</FormLabel>
                        <FormDescription>Spin up a per-user/container instance.</FormDescription>
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value}
                          onCheckedChange={(val) => {
                            field.onChange(val);
                            if (val) {
                              setValue("docker", {
                                image: "",
                                compose: "",
                                hashDomain: false,
                                lifetime: 3600,
                                envs: "",
                                maxMemory: 0,
                                maxCPU: 0,
                              }, { shouldValidate: true, shouldDirty: true });
                            } else {
                              setValue("docker", undefined, { shouldValidate: true, shouldDirty: true });
                            }
                          }}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

        {instanced && (
          <div className="space-y-4 rounded-lg border p-4">
            <h3 className="text-base font-semibold flex"><Container className="mr-2" size={24} /> Docker Config</h3>

            <FormField
              control={form.control}
              name="docker.image"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Image</FormLabel>
                  <FormControl>
                    <Input placeholder="registry/image:tag" {...field} />
                  </FormControl>
                  <FormDescription>Docker image to run.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="docker.compose"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Compose (YAML)</FormLabel>
                  <FormControl>
                    <div className="rounded-md border" style={{ maxHeight: "400px", overflow: "auto" }}>
                      <CodeMirror
                        value={field.value ?? ""}
                        minHeight="200px"
                        maxHeight="400px"
                        extensions={yamlExtensions}
                        theme={oneDark}
                        onChange={(val) => field.onChange(val)}
                        style={{ width: "100%", minWidth: 0 }}
                        basicSetup={{
                          lineNumbers: true,
                          highlightActiveLine: true,
                          highlightActiveLineGutter: true,
                          foldGutter: true,
                        }}
                      />
                    </div>
                  </FormControl>
                  <FormDescription>Optional docker-compose override.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="docker.lifetime"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Lifetime (seconds)</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      min={0}
                      step={60}
                      value={field.value ?? 0}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormDescription>Duration before instance cleanup.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="docker.hashDomain"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <FormLabel>Hash domain</FormLabel>
                    <FormDescription>Unique subdomain per instance.</FormDescription>
                  </div>
                  <FormControl>
                    <Switch checked={!!field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="docker.envs"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Environment variables</FormLabel>
                  <FormControl>
                    <Textarea placeholder="KEY=VALUE" {...field} />
                  </FormControl>
                  <FormDescription>Line-separated env vars.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="docker.maxMemory"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Max Memory (MB) <MemoryStick size={24} /> </FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min={0}
                        step={128}
                        value={field.value ?? 0}
                        onChange={(e) => field.onChange(Number(e.target.value))}
                      />
                    </FormControl>
                    <FormDescription>0 = default/unlimited</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="docker.maxCPU"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Max CPU <Cpu size={24} /></FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min={0}
                        step={0.1}
                        value={field.value ?? 0}
                        onChange={(e) => field.onChange(Number(e.target.value))}
                      />
                    </FormControl>
                    <FormDescription>0 = default/unlimited</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </div>
        )}

        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}
